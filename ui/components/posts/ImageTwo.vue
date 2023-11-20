<template>
  <div class="p-1">
    <div class="flex h-80 w-full">
      <img class="rounded-lg  w-1/2 object-cover"
        :src="`/streaming-api/post/${postId}/image/${imageIds[0]}`"
        id="1"
        v-on:click="changeClickedIndex(0)"
      />
      <img class="w-1/2 rounded-lg object-cover ml-0.5"
        :src="`/streaming-api/post/${postId}/image/${imageIds[1]}`"
        id="2"
        v-on:click="changeClickedIndex(1)"
        v-on:closeModal="closeModal"
      />
    </div>
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
  showModal.value = false
}
</script>