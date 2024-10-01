<script>
export default {

  props: ["photos",],

  methods: {
    getBase64Image(imageData) {
      // Controlla se imageData è valido
      if (!imageData || imageData.length === 0) {
        console.error("Nessun dato immagine valido fornito.");
        return ''; // Restituisci una stringa vuota se non ci sono dati validi
      }

      try {
        // Creiamo un Uint8Array dall'array di byte
        const byteArray = new Uint8Array(imageData);
        console.log("Byte Array:", byteArray);
        // Convertiamo i byte in una stringa binaria
        const binaryString = byteArray.reduce((data, byte) => data + String.fromCharCode(byte), '');
        console.log("Binary String:", binaryString);

        // Convertiamo la stringa binaria in base64
        const base64String = btoa(binaryString);
        console.log("Base64 String:", base64String);

        // Restituiamo la stringa formattata per l'attributo src dell'immagine
        return 'data:image/jpeg;base64,' + base64String; // Modifica il prefisso se necessario
      } catch (error) {
        console.error("Errore nella conversione dei dati immagine:", error);
        return ''; // Restituisci una stringa vuota in caso di errore
      }
    }
  },

}

</script>

<template>
  <div>
    <div v-for="photo in photos" :key="photo.photo_Id" class="photo">
      <h2>{{ photo.owner.nickname }}</h2>
      <img src= "data:image/png;base64,photo.image" />
      <button @click="toggleLike">
        <i :class="['bi', 'bi-hand-thumbs-up', { 'text-danger': photo.liked, 'text-muted': !photo.liked }]"></i>
      </button>
      <p>Date: {{ photo.date }}</p>
      <p>Likes: {{ photo.nLikes }}</p>
    </div>
  </div>
</template>

<style scoped>

</style>