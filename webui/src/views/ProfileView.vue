<script>
import Photo from "../components/Photo.vue";

export default {
  components: {Photo},
  data() {
    return {
      errormsg: null,
      nickname: "",
      followers: [],
      following: [],
      photos: [],
      nFollowers: 0,
      nFollowing: 0,
      isFollowing: false,
      existsBan: false,
      requester: 0,
      photoOpen: false,
    };
  },

  props: ["id"],

  mounted() {

  },

  created() {

    this.getProfile(this.id);
  },

  methods: {
    async getProfile(id) {
      try {
        let response = await this.$axios.get(`/users/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
        });

        // Assegna i dati dalla risposta alle proprietà del componente
        this.nickname = response.data.nickname;
        this.followers = response.data.followers;
        this.following = response.data.following;
        this.photos = response.data.photos;
        this.nFollowers = response.data.nFollowers;
        this.nFollowing = response.data.nFollowing;
        this.isFollowing = response.data.isFollowing;
        this.existsBan = response.data.existsBan
        this.requester = localStorage.getItem('token')

      } catch (e) {
        this.errormsg = e.toString();

      }
    },
  },
};
</script>

<template>
  <div>
    <div class="profile-bar">
      <div class="left-section">
        <span>{{ nickname }}</span>
      </div>
      <div class="right-section">
        <p>Followers: {{ nFollowers }}</p>
        <p style="margin-top: 8px;">Following: {{ nFollowing }}</p>
      </div>
    </div>

    <div v-if="Number(requester) === id" class="button-section">
      <button @click="toggleFollow" class="button" :style="{ backgroundColor: isFollowing ? 'orange' : 'blue' }">
        {{ isFollowing? "UnFollow" : "Follow" }}
      </button>
      <button @click="toggleBlock" class="button" :style="{ backgroundColor: existsBan ? 'green' : 'red'}">
        {{ existsBan? "Unblock" : "Block" }}
      </button>
    </div>
  </div>


  <div class="flex-container">
    <div v-for="photo in photos" :key="photo.photo_Id"  class="mini-card" >
      <img :src="'data:image/png;base64, ' + photo.image" />
    </div>
  </div>



</template>

<style>
.profile-bar {
  display: flex; /* Usa flexbox per una migliore disposizione */
  justify-content: space-between; /* Spazio tra le sezioni */
  align-items: flex-start; /* Allinea gli elementi in alto */
  padding: 20px; /* Spazio interno per il contenitore */
}

.left-section {
  text-align: left;
  font-size: 34px; /* Riduci la grandezza del font per il nickname */
  font-weight: bold;
  margin-left: 20px;
}

.right-section {
  text-align: right;
  font-size: 16px; /* Dimensione del font per i follower e following */
  font-weight: bold;
  margin-right: 20px;
}

.button-section {
  text-align: right;
  font-size: 16px; /* Riduci la grandezza del font per il nickname */
  font-weight: bold;
  margin-right: 20px;
  gap: 20px;
}
.button{
  margin-right: 20px;
  width: 250px;

}

.flex-container {
  display: flex;
  flex-wrap: wrap; /* Permette di andare a capo quando non c'è più spazio */
  margin: 10px;                 /* Rimuove margini esterni */
  padding: 0;                /* Rimuove padding interno */
}

.mini-card {
  width: 20%;
  height: 300px;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  margin: 0;
  padding:0;
}

.mini-card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.mini-card:hover {
  transform: translateY(-10px); /* Solleva l'elemento di 5px verso l'alto */
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2); /* Aggiunge un'ombra per l'effetto di elevazione */
}

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


</style>