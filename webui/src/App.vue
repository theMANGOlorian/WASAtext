<template>
	<header class="navbar navbar-dark sticky-top custom-purple flex-md-nowrap p-0 shadow d-flex justify-content-start align-items-center">
	  <a class="navbar-brand custom-title col-md-3 col-lg-2 me-0 px-3 fs-4" href="#">WasaText</a>
	  <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
		<span class="navbar-toggler-icon"></span>
	  </button>
    <template v-if="$route.path !== '/session'">
      <GoToButton :text="'Chats'" :path="'chats'"></GoToButton>
      <GoToButton :text="'Profile'" :path="'settings'"></GoToButton>
    </template>
	</header>
  
	<RouterView />
</template>

<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<script>
export default {}
</script>

<style>

html, body {
    margin: 0;
    padding: 0;
    height: 100%; /* Assicura che l'intera pagina prenda tutta l'altezza */
    overflow: hidden; /* Impedisce lo scrolling dell'intera pagina */
}

</style>
  