<script>
export default {

  props: {
    photos: {
      type: Array,
      required: true
    },
    id: {
      type: Number,
      required: true
    },
  },

  data() {
    return {
      errormsg: null,
    }
  },



  methods: {


    // Funzione per cambiare lo stato dell'oggetto
    async toggleLike(id) {
      // Trova l'oggetto da modificare in base all'ID
      const p = this.photos.find(photo => photo.photo_Id === id);
      if (p) {
        // Inverte il valore di isTrue
        p.liked = !p.liked;
        if (p.liked) {
          try {
            await this.$axios.put("/users/" + this.id + "/photos/" + p.photo_Id + "/likes/" + localStorage.getItem('token'), {
              headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
          });
            p.nLikes = p.nLikes +1
          } catch (e) {
            p.liked = !p.liked;
            this.errormsg = e.toString();
          }
        } else {
          try {
            await this.$axios.delete("/users/" +  this.id  + "/photos/" + p.photo_Id + "/likes/" + localStorage.getItem('token'), {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              },
            });
            p.nLikes = p.nLikes -1
          } catch (e) {
            p.liked = !p.liked;
            this.errormsg = e.toString();
          }
        }
      }
    },
  },
}

</script>

<template>
  <div>
    <div v-for="photo in photos" :key="photo.photo_Id" class="photo">
      <h2>{{ photo.owner.nickname }} </h2>
      <img src= "data:image/png;base64,photo.image" />
      <div class="info-container">
        <span class="like-text"> Comment </span>
        <button class="like-button"></button>
        <span class="like-text"> Put a Like  </span>
        <button @click="toggleLike(photo.photo_Id)"
                class="like-button"
                :class="photo.liked ? 'text-primary': '' "    >
        </button>
        <span> {{ photo.liked}}</span>
        <span class="like-text"> Likes:  {{ photo.nLikes }} </span>
      </div>
    </div>
  </div>
</template>

<style scoped>

.info-container {
  display: flex;               /* Flexbox per allineare orizzontalmente */
  align-items: center;         /* Allineamento verticale */
}


.like-text {
  font-size: 15px;            /* Dimensione del testo più piccola */
  margin-left: 20px;           /* Margine a sinistra per separare dall'icona */
  margin-right: 5px;
}


.like-button {


  background-color: #f0f0f0; /* Sfondo chiaro */
  border: none;               /* Nessun bordo */
  border-radius: 50%;         /* Bottone circolare */
  padding: 10px;              /* Spazio interno */
  cursor: pointer;            /* Puntatore a mano */
  width: 20px;                /* Larghezza del bottone */
  height: 20px;               /* Altezza del bottone */
  display: flex;
  align-items: center;
  justify-content: center;    /* Icona centrata */
  transition: background-color 0.3s ease; /* Animazione smooth */
}

.like-button:hover {
  background-color: #1877f2; /* Colore di sfondo quando il bottone è hoverato */
}

.like-button {
  font-size: 30x; /* Dimensione dell'icona */

}

.text-primary {
  background-color: #1877f2 !important; /* Colore blu per il like attivo, simile a Facebook */
}

.like-button:hover i {
  color: #6c757d;  /* Cambia colore al passaggio del mouse */
}


</style>