<template>
  <div class="p-1">
    <img :src="`/streaming-api/post/${postId}/image/${imageIds[0]}`"
      class="rounded-t-lg min-w-full object-scale-down h-90"
      v-on:click="changeClickedIndex(0)"
    />
  </div>
  <Modal v-on:closeModal="closeModal" :postId="postId" :imageIds="imageIds" :initialIndex="0" v-if="showModal"/>
</template>

<script setup lang="ts">
const props = defineProps({
	imageIds: { type: Array<string>, required: true },
  postId: { type: String, required: true},
})

const showModal = ref(false)
const clickedIndex = ref(0)
const changeClickedIndex = function(index: number) {
  showModal.value = true
  clickedIndex.value = index
}

const closeModal = function() {
  console.log('showModal to false')
  showModal.value = false
}
</script>
