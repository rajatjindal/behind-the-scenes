<template>
	<div class="mx-auto w-2/3 px-1 lg:px-0 mt-10 grid grid-cols-2 gap-4">
		<div v-for="postId in postkeys" :key="postId" class="col-span-1">
			<Post :postId="postId" />
		</div>
	</div>
</template>
  
<script setup lang="ts">
import { myfetch } from "@/sdk/base/myfetch";

const postkeys = ref<string>()
const loading = ref<Boolean>(false)
const latency = ref<Number>(0)

onBeforeMount(async () => {
	loading.value = true
	const startTime = new Date();

	postkeys.value = await myfetch('/posts', { method: 'GET'})

	loading.value = false
	const endTime = new Date();
	latency.value = Math.abs(endTime.getMilliseconds() - startTime.getMilliseconds());
})
</script>