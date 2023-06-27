<template>
  <div class="">
    <img :src="`https://fermyon-bts.usingspin.com/api/post/${postId}/image/${imageIds[0]}`"
      class="min-w-full h-80 object-cover"
    />
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
	imageIds: { type: Array<string>, required: true },
  postId: { type: String, required: true},
})
</script>
