<script>
import Photo from "../components/Photo.vue";

export default {
  components: {Photo},
  data() {
    return {
      errormsg: null,
      nicknameProfileOwner: "",
      followers: [],
      following: [],
      photos: [],
      nFollowers: 0,
      nFollowing: 0,
      isFollowing: false,
      existsBan: false,
      requester: localStorage.getItem('token'),
      modalPhoto: false,
      selectedPhoto: {},
      id: this.$route.params.id,
    };
  },


  created() {

    this.getProfile(this.id)
  },
  watch: {
    // Quando l'ID nella route cambia, carica di nuovo i dati
    '$route.params.id': function(newId) {
      this.getProfile(Number(newId)); // Carica i dati con il nuovo ID
    }
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
        this.nicknameProfileOwner = response.data.nickname;
        this.followers = response.data.followers;
        this.following = response.data.following;
        this.photos = response.data.photos;
        this.nFollowers = response.data.nFollowers;
        this.nFollowing = response.data.nFollowing;
        this.isFollowing = response.data.isFollowing;
        this.existsBan = response.data.existsBan;
        this.requester = localStorage.getItem('token');
        this.errormsg = null;
        this.modalPhoto = false;
        this.selectedPhoto = [];

      } catch (e) {
        this.errormsg = e.toString();

      }
    },
    async openPhoto(photo) {
      try {
        const response = await this.$axios.get(`/users/${this.requester}/photos/${photo.photo_Id}`, {
          headers: {
            "Authorization": `Bearer ${localStorage.getItem('token')}`,
          }
        });
        this.selectedPhoto = response.data;
        this.modalPhoto = true;
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    closePhoto(){
      this.modalPhoto=false;
      this.selectedPhoto= {};
    },
    async toggleFollow() {
      if (this.isFollowing) {
        try {
          await this.$axios.delete("/users/" + this.requester + "/followers/" + this.id, {
            headers: {
              "Authorization": `Bearer ${localStorage.getItem('token')}`
            },
          });
          await this.getProfile(this.id)
        } catch (e) {
          console.log(e.toString());
        }
      } else {
        try {
          await this.$axios.put("/users/" + this.requester + "/followers/" + this.id, {
            headers: {
              "Authorization": `Bearer ${localStorage.getItem('token')}`
            },
          });
          await this.getProfile(this.id)
        } catch (e) {
          console.log(e.toString());
        }
      }
    },
    async toggleBan(){
      if (this.existsBan) {
        try {
          await this.$axios.delete("/users/" + this.requester + "/banned_users/" + this.id, {
            headers: {
              "Authorization": `Bearer ${localStorage.getItem('token')}`
            },
          });
          await this.getProfile(this.id)
        } catch (e) {
          console.log(e.toString());
        }
      } else {
        try {
          await this.$axios.put("/users/" + this.requester + "/banned_users/" + this.id, {
            headers: {
              "Authorization": `Bearer ${localStorage.getItem('token')}`
            },
          });
          await this.getProfile(this.id)
        } catch (e) {
          console.log(e.toString());
        }
      }
    },
  },
};

</script>

<template>

  <div v-if="!modalPhoto">


    <div>
      <div class="profile-bar">
        <div class="left-section">
          <span>{{ nicknameProfileOwner }}</span>
        </div>
        <div class="right-section">
          <p>Followers: {{ nFollowers }}</p>
          <p style="margin-top: 8px;">Following: {{ nFollowing }}</p>
        </div>
      </div>

      <div v-if="Number(requester) !== Number(id)" class="button-section">
        <button v-if= "!existsBan " @click="toggleFollow" class="button" :style="{ backgroundColor: isFollowing ? 'orange' : 'blue' }">
          {{ isFollowing? "UnFollow" : "Follow" }}
        </button>
        <button @click="toggleBan" class="button" :style="{ backgroundColor: existsBan ? 'green' : 'red'}">
          {{ existsBan? "Unblock" : "Block" }}
        </button>
      </div>
    </div>


    <div class="flex-container">
      <div v-for="photo in photos" :key="photo.photo_Id"  class="mini-card" @click="openPhoto(photo)">
        <img :src="'data:image/png;base64, ' + photo.image" />
      </div>
    </div>

  </div>

  <div v-if="modalPhoto"  class="overlay-background" @click.self="closePhoto">
    <photo :photos="[selectedPhoto]"
           @photoDeleted="getProfile(id)"
    >
    </photo>

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

.overlay-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.8); /* Trasparenza per lo sfondo */
  display: flex;
  align-items: center;
  justify-content: center; /* Centra la modale */
  z-index: 1000;
  overflow: hidden;

}
photo{
  color:white;
}



</style>