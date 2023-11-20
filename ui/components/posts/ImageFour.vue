<template>
  <div class="p-1">
    <div class="flex h-80 w-full">
      <div class="flex w-1/2 flex-col ml-0.5">
        <div class="h-1/2">
          <img
            class="h-full w-full rounded-lg "
            :src="`/streaming-api/post/${postId}/image/${imageIds[0]}`"
            v-on:click="changeClickedIndex(0)"
            id="1"
          />
        </div>
        <div
          class="h-1/2 rounded-bl-xl mt-0.5"
        >
          <img
            class="rounded-lg h-full w-full"
            :src="`/streaming-api/post/${postId}/image/${imageIds[1]}`"
            v-on:click="changeClickedIndex(1)"
            id="2"
          />
        </div>
      </div>

      <div class="flex w-1/2 flex-col ml-0.5">
        <div class="h-1/2 rounded-tr-xl">
          <img
            class="rounded-lg h-full w-full"
            :src="`/streaming-api/post/${postId}/image/${imageIds[2]}`"
            v-on:click="changeClickedIndex(2)"
            id="3"
          />
        </div>
        <div
          class="h-1/2 rounded-br-xl mt-0.5"
        >
          <img
            class="rounded-lg h-full w-full"
            :src="`/streaming-api/post/${postId}/image/${imageIds[3]}`"
            v-on:click="changeClickedIndex(3)"
            id="4"
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
