<template>
  <div class="h-64 w-full mt-5 p-0.5">
    <div class="flex h-full w-full">
      <img class="w-1/2 border border-gray-300 border-2 rounded-l-lg overflow-hidden"
        :src="'data:image/png;base64,' + imageIds[0]"
        id="1"
      />
      <img class="w-1/2 border border-gray-300 border-2 rounded-r-lg ml-0.5 overflow-hidden"
      :src="'data:image/png;base64,' + imageIds[0]"
        id="2"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
	imageIds: { type: Array<string>, required: true },
})
</script>