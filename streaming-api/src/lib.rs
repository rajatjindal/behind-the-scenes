use spin_sdk::http_component;

use anyhow::{anyhow, Result};
use futures::SinkExt;
use futures::StreamExt;
use serde::{Deserialize, Serialize};
use spin_sdk::{
    http::{
        self, Headers, IncomingRequest, IncomingResponse, Method, OutgoingRequest,
        OutgoingResponse, ResponseOutparam, Scheme,
    },
    key_value::Store,
    variables,
};
use url::Url;

use flate2::{write::GzEncoder, Compression};
use std::io::Write;

/// Send an HTTP request and return the response.
#[http_component]
async fn send_outbound(req: IncomingRequest, res: ResponseOutparam) {
    get_and_stream_imagefile(req, res).await.unwrap()
}

#[derive(Clone, Debug, Default, Deserialize, Serialize)]
#[serde(deny_unknown_fields, rename_all = "camelCase")]
pub struct Post {
    pub msg: String,
    pub timestamp: String,
    pub image_ids: Vec<String>,
    pub image_map: std::collections::HashMap<String, String>,
    pub approved: bool,
    pub grapes: i64,
    pub hearts: i64,
}

async fn get_and_stream_imagefile(req: IncomingRequest, res: ResponseOutparam) -> Result<()> {
    let (post_id, image_id) =
        get_post_id_and_image_id_from_url(req.path_with_query().unwrap_or_default());

    let store = Store::open_default()?;
    let raw_post = store.get(format!("post:{post_id}").as_str())?.unwrap();

    let post: Post = serde_json::from_slice(&raw_post)?;
    let url = post.image_map.get(image_id.as_str()).unwrap();

    let token = variables::get("slack_token").unwrap();
    let url = Url::parse(url).unwrap();
    let outgoing_request = OutgoingRequest::new(
        &Method::Get,
        Some(url.path()),
        Some(&match url.scheme() {
            "http" => Scheme::Http,
            "https" => Scheme::Https,
            scheme => Scheme::Other(scheme.into()),
        }),
        Some(url.authority()),
        &Headers::new(&[(
            "authorization".to_string(),
            format!("Bearer {token}").as_bytes().to_vec(),
        )]),
    );

    let response = http::send::<_, IncomingResponse>(outgoing_request).await?;

    let status = response.status();
    if status != 200 {
        return Err(anyhow!(format!(
            "failed to fetch image from slack. expected 200, got status code {}",
            status
        )));
    }

    let mut stream = response.take_body_stream();

    let content_type = match response.headers().get("content-type").first() {
        Some(content_type) => content_type.to_owned(),
        None => b"image/jpeg".to_vec(),
    };

    let out_response = OutgoingResponse::new(
        status,
        &Headers::new(&[
            ("content-type".to_string(), content_type),
            ("content-encoding".to_string(), b"gzip".to_vec()),
        ]),
    );

    let mut body = out_response.take_body();
    res.set(out_response);

    let mut encoder = GzEncoder::new(Vec::new(), Compression::fast());
    while let Some(chunk) = stream.next().await {
        // let chunk = chunk?;
        encoder.write_all(&chunk?).unwrap();
        body.send(encoder.get_ref().to_vec()).await?;
        encoder.get_mut().clear()
    }

    Ok(())
}

fn get_post_id_and_image_id_from_url(path_and_query: String) -> (String, String) {
    let path = path_and_query.replace("/streaming-api/post/", "");
    let parts: Vec<&str> = path.split("/").collect();
    let post_id = parts[0];
    let image_id = parts[2];

    return (post_id.to_string(), image_id.to_string());
}
