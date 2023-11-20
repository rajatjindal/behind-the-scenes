<template>
  <div class="p-1">
    <div class="flex h-80 w-full">
      <div class="w-1/2 rounded-lg object-cover">
        <img 
          :src="`/streaming-api/post/${postId}/image/${imageIds[0]}`" 
          id="1"
          v-on:click="changeClickedIndex(0)"
          class="h-full w-full" />
      </div>
      <div class="flex w-1/2 flex-col ml-0.5">
        <div class="h-1/2 rounded-lg object-cover">
          <img
            class="h-full w-full"
            :src="`/streaming-api/post/${postId}/image/${imageIds[1]}`"
            v-on:click="changeClickedIndex(1)"
            id="2"
          />
        </div>
        <div
          class="h-1/2 rounded-lg object-cover py-0.5"
        >
          <img
            class="h-full w-full"
            v-on:click="changeClickedIndex(2)"
            :src="`/streaming-api/post/${postId}/image/${imageIds[2]}`"
            id="3"
          />
        </div>
      </div>
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