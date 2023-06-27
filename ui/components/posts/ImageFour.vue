<template>
  <div class="w-full mt-5 p-0.5 pr-2 lg:pr-0">
    <div class="flex h-full w-full">
      <div class="flex w-1/2 flex-col ml-0.5">
        <div class="h-1/2 border border-gray-300 overflow-hidden rounded-tl-xl">
          <img
            class="h-full w-full"
            :src="'data:image/png;base64,' + imageIds[0]"
            id="1"
          />
        </div>
        <div
          class="h-1/2 border border-gray-300 overflow-hidden rounded-bl-xl mt-0.5"
        >
          <img
            class="h-full w-full"
            :src="'data:image/png;base64,' + imageIds[0]"
            id="2"
          />
        </div>
      </div>

      <div class="flex w-1/2 flex-col ml-0.5">
        <div class="h-1/2 border border-gray-300 overflow-hidden rounded-tr-xl">
          <img
            class="h-full w-full"
            :src="'data:image/png;base64,' + imageIds[0]"
            id="3"
          />
        </div>
        <div
          class="h-1/2 border border-gray-300 overflow-hidden rounded-br-xl mt-0.5"
        >
          <img
            class="h-full w-full"
            :src="'data:image/png;base64,' + imageIds[0]"
            id="4"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
	imageIds: { type: Array<string>, required: true },
})
</script>
