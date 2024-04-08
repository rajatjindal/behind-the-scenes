<template>
	<div class="mt-1" v-if="runsOn === 'fermyon'">
		<a href="https://www.fermyon.com/cloud"><img src="https://www.fermyon.com/static/image/fermyon-badge.png" /></a>
	</div>
	<div class="mt-2" v-else-if="runsOn === 'spin'">
		<a href="https://www.fermyon.com/cloud"><img class="h-12" src="https://github.com/fermyon/spin/raw/main/docs/static/image/logo-dark.png" /></a>
	</div>
	<div class="w-40 h-16" v-else-if="runsOn === 'spinkube'">
		<a href="https://spinkube.dev">
			<SpinKube class="w-40 h-full"/>
		</a>
	</div>
</template>

<script setup lang="ts">
import { myfetch } from "@/sdk/base/myfetch";
const runsOn = ref<string>("fermyon")

const updateRunsOn = async function() {
	runsOn.value = await myfetch(`https://fermyon-bts.usingspin.com/api/runs-on`, { method: 'GET'})
}

onBeforeMount(async () => {
	await updateRunsOn()
})
</script>