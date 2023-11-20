<template>
	<div class="w-full mx-auto w-full px-1 lg:px-0 border rounded-t-xl border-b-0" v-if="post">
		<ImageOne v-if="post && post.imageIds.length === 1" :postId="post.timestamp" :imageIds="post.imageIds" />
		<ImageTwo v-if="post && post.imageIds.length === 2" :postId="post.timestamp" :imageIds="post.imageIds" />
		<ImageThree v-if="post && post.imageIds.length === 3" :postId="post.timestamp" :imageIds="post.imageIds" />
		<ImageFour v-if="post && post.imageIds.length === 4" :postId="post.timestamp" :imageIds="post.imageIds" />
		<div class="border-b px-4 py-2 flex">
			<div class="flex"><Grapes class="w-8 h-8" fill="indigo" v-on:click="incrementGrapes"/><span class="ml-2 text-gray-400 flex items-center justify-between text-lg">{{ post.grapes || 0}}</span></div>
			<div class="flex ml-5"><Heart class="w-8 h-8" fill="red" v-on:click="incrementHearts"/><span class="ml-2 text-gray-400 flex items-center justify-between text-lg">{{ post.hearts || 0}}</span></div>
		</div>
	</div>
</template>
  
<script setup lang="ts">
import { myfetch } from "@/sdk/base/myfetch";

const props = defineProps({
	postId: { type: String, required: true },
})

declare interface Post {
	msg:       string;
	timestamp: string;
	imageIds:    string[];
	grapes: number,
	hearts: number
}


const post = ref<Post>()
const loading = ref<Boolean>(false)
const latency = ref<Number>(0)

onBeforeMount(async () => {
	await reloadPost()
})

const reloadPost = async function() {
	loading.value = true
	const startTime = new Date();

	post.value = await myfetch(`/post/${props.postId}`, { method: 'GET'})

	loading.value = false
	const endTime = new Date();
	latency.value = Math.abs(endTime.getMilliseconds() - startTime.getMilliseconds());
}

const withoutSignoff = function(input: string): string {
	return input.replaceAll('/signoff', '')
}

const incrementGrapes = async function() {
	await myfetch(`/post/${props.postId}/grapes`, { method: 'POST'})
	await reloadPost()
}

const incrementHearts = async function() {
	await myfetch(`/post/${props.postId}/hearts`, { method: 'POST'})
	await reloadPost()
}

</script>