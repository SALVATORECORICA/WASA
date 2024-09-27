<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
  data(){
    return {
      logged: false,
      searchValue: ""
    }
  },
  created() {
    if (!localStorage.getItem('notFirstStart')) {
      localStorage.clear()
      localStorage.setItem('notFirstStart', true)
      // console.log("first start")
    }
  },
  mounted(){
      if (!localStorage.getItem('token')){
        this.$router.replace("/login")
      }else{
        this.logged = true
      }
    },
}


</script>

<template>
<div>
	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Wasagram</a>
    <form v-if="!logged" class="d-flex ms-auto me-3 input-group" aria-label="Insert Nickname" style="max-width: 300px;">
      <input class="form-control" type="search" placeholder="Insert Nickname"  style ="width:150px;">
      <button class="btn btn-outline-light" type="submit">Enter</button>
    </form>
	</header>

	<div v-if="logged" class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link1" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Menu item 1
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link2" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Menu item 2
							</RouterLink>
						</li>
					</ul>

					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Secondary menu</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
								Item 1
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</div>
</template>

<style>
.btn {
  margin: 0; /* Rimuove margini */
  padding: 0.375rem 0.75rem; /* Regola il padding del pulsante, se necessario */
}

.btn-outline-light:hover {
  background-color: red; /* Cambia lo sfondo in rosso al passaggio del mouse */
  color: #fff; /* Il colore del testo rimane bianco */
  border-color: red; /* Cambia il colore del bordo in rosso al passaggio del mouse */
}
</style>
