<script>

export default{
  data(){
    return{
    }
  },
  props: ["modalSearchOn",],

  methods : {
    closeModalSearch(){
      this.$emit("closeModalSearch")
    },
  },
}
</script>



<template>

  <div  v-if="modalSearchOn" class="overlay-background" @click.self="closeModalSearch">
    <div class="overlay">
      <span class="label">aaaaaaa</span>
      <span class="label">aaaaaaa</span>

    </div>
  </div>

</template>


<style scoped>
.overlay-background{
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.8); /* Trasparenza per lo sfondo */
  display: flex;
  align-items: center;
  justify-content: center; /* Centra la modale */
  z-index: 1000;
  overflow: hidden;

}


.overlay{
  background-color: white;
  padding: 20px;
  width: 400px; /* Larghezza fissa */
  height: 350px; /* Altezza fissa */
  overflow-y: auto; /* Abilita lo scroll solo per il contenuto */
  border-radius: 10px;
  flex-direction: column;
}
.label{
  background-color: lightgray; /* Colore di sfondo */
  color: black; /* Colore del testo */
  padding: 10px 15px; /* Padding interno per la grandezza */
  border-radius: 5px; /* Bordo arrotondato */
  cursor: pointer; /* Cursore a forma di mano al passaggio del mouse */
  transition: background-color 0.3s; /* Transizione per il cambio di colore */
  display: block; /* Imposta il display come block */
  width: 100%;
  border: 2px solid black;
  margin-bottom: 10px;
}

.label:hover {
  background-color: #0056b3; /* Colore di sfondo al passaggio del mouse */
}

</style>