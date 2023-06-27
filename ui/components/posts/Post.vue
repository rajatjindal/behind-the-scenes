<template>
	<div class="w-full mx-auto w-full px-1 lg:px-0 border rounded-t-xl" v-if="post">
		<div class="border-b px-4 py-2">
			<span class="text-gray-600 italic">{{ withoutSignoff(post.msg) }}</span>
		</div>
			<ImageOne v-if="post && post.imageIds.length === 1" :postId="post.timestamp" :imageIds="post.imageIds" />
			<ImageTwo v-if="post && post.imageIds.length === 2" :postId="post.timestamp" :imageIds="post.imageIds" />
			<ImageThree v-if="post && post.imageIds.length === 3" :postId="post.timestamp" :imageIds="post.imageIds" />
			<ImageFour v-if="post && post.imageIds.length === 4" :postId="post.timestamp" :imageIds="post.imageIds" />
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
}


const post = ref<Post>()
const loading = ref<Boolean>(false)
const latency = ref<Number>(0)

onBeforeMount(async () => {
	loading.value = true
	const startTime = new Date();

	post.value = await myfetch(`/post/${props.postId}`, { method: 'GET'})

	loading.value = false
	const endTime = new Date();
	latency.value = Math.abs(endTime.getMilliseconds() - startTime.getMilliseconds());
})

const withoutSignoff = function(input: string): string {
	return input.replaceAll('/signoff', '')
}
</script>