<script>




export default {
  data() {
    return {
      er: null,
      nickname : ""
    }
  },
  mounted() {
    if (localStorage.getItem('token')){
      this.$router.replace("/home")
    }
    console.log("LoginView mounted!");
  },
  methods:{
    async login(){
      this.errormsg ="";
      try{
        let response = await this.$axios.post("/session", {nickname: this.nickname.trim()});
        localStorage.setItem("token", response.data.id)
        this.$router.replace("/home")
        this.$emit("endLogin",response.data.id, this.nickname)
      }
      catch (err){
        this.er = err.toString()
        console.log(this.er)
      }
    }
  }
}
</script>


<template>
  <div class="login">
    <form  class="d-flex ms-auto me-3 input-group" aria-label="Insert Nickname" style="max-width: 300px;" @submit.prevent="login">
      <input
          class="form-control"
          type="search"
          placeholder="Insert Nickname"
          style ="width:150px;" v-model="nickname"
          maxlength="16"
          minlength="3">
      <button class="btn btn-outline-light" type="submit" >Enter</button>
    </form>
  </div>
</template>

<style>
.login {
  background-image: url("../../public/Blog_SimpleStrategies_030421-1.jpg"); /* Percorso all'immagine */
  background-size: cover; /* Assicura che l'immagine riempia l'elemento */
  height: 100vh; /* Altezza della viewport */
  width: 100%;
  position: relative; /* Consente l'uso di posizioni assolute nei contenuti figli, se necessario */
  margin: 0;
  padding: 0;
  overflow: hidden;
}
</style>