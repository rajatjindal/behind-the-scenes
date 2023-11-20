<template>
  <div class="fixed z-10 inset-0 overflow-y-auto"
       aria-labelledby="modal-title"
       role="dialog"
       aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-900 bg-opacity-95 transition-opacity"
           aria-hidden="true"
           v-on:click="closeModal"></div>

      <div class="fixed bottom-8 mx-auto sm:top-4 sm:right-4 z-20 text-white" v-on:click="closeModal">
        <div class="text-gray-100 bg-gray-900 rounded-full font-bold w-16 h-16 flex items-center justify-center">
          <svg width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17.25 6.75L6.75 17.25"></path>
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6.75 6.75L17.25 17.25"></path>
          </svg>
        </div>
      </div>

      <div v-if="currentIndex > 0"
           class="fixed left-4 bottom-4 sm:bottom-1/2">
        <div class="text-gray-100 bg-gray-900 rounded-full font-bold w-16 h-16 flex items-center justify-center"
             v-on:click="prev">
          <svg width="24"
               height="24"
               fill="none"
               viewBox="0 0 24 24">
            <path stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M10.25 6.75L4.75 12L10.25 17.25" />
            <path stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M19.25 12H5" />
          </svg>
        </div>
      </div>

      <div v-if="currentIndex < imageIds.length - 1"
           class="fixed right-4 bottom-4 sm:bottom-1/2 z-20">
        <div class="text-gray-100 bg-gray-900 rounded-full py-5 font-bold w-16 h-16 flex items-center justify-center"
             v-on:click="next">
          <svg width="24"
               height="24"
               fill="none"
               viewBox="0 0 24 24">
            <path stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M13.75 6.75L19.25 12L13.75 17.25" />
            <path stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.5"
                  d="M19 12H4.75" />
          </svg>
        </div>
      </div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen"
            aria-hidden="true">&#8203;</span>
      <div class="inline-block my-auto rounded-lg text-left shadow-xl transform transition-all sm:align-middle sm:my-8 sm:max-w-lg sm:w-full">
        <div class="overflow-hidden rounded-2xl">
          <img :src="`/streaming-api/post/${postId}/image/${imageIds[currentIndex]}`" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits(['closeModal'])

const props = defineProps({
	imageIds: { type: Array<string>, required: true },
  postId: { type: String, required: true},
  initialIndex: { type: Number, required: true}
})

const currentIndex = ref(props.initialIndex)
const prev = function() {
  currentIndex.value = currentIndex.value - 1
}

const next = function() {
  currentIndex.value = currentIndex.value + 1
}

const closeModal = function() {
  console.log('close called')
  emit('closeModal', '')
}
</script>