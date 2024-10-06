<script>
export default {

  props: {
    photos: {
    },
    id: {   // id of the logged User
      type: Number,
      required: true
    },
    nickname: {  // nickname of the logged User

    }
  },

  data() {
    return {
      errormsg: null,
      commentModal:false,
      openedPhoto: null,
      newComment:"",
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
    openComments(photo){
      this.openedPhoto=photo;
      this.commentModal=true;
    },
    closeComments(){
      this.openedPhoto=null;
      this.commentModal=false;
      this.newComment="";
    },
    async addComment(){
      try {
        await this.$axios.post("/users/" +  this.id  + "/photos/" + this.openedPhoto.photo_Id + "/comments" , {
            "comment": this.newComment,
          }, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
        });
      } catch(e) {
      }
    }
  },
}

</script>

<template>
  <div v-if="commentModal" @click.self="closeComments" class="overlay-background">
    <div class="overlay">
      <div class="comments-container">
        <div v-for="comment in openedPhoto.comments">
          <span> {{ comment.user}}</span>
          <span> {{comment.comment}}</span>
          <span> Delete Comment</span>
        </div>
      </div>
      <div class="fixed-bottom">
        <input placeholder="Insert your comment" v-model="newComment">
        <button @click="addComment"> Add comment</button>
      </div>
    </div>
  </div>

  <div v-if="!commentModal">
    <div v-for="photo in photos" :key= " photo.photo_Id" class="photo" @click="openComments(photo)">
      <h2>{{ photo.owner.nickname }} </h2>
      <img :src= "'data:image/png;base64, ' + photo.image" />
      <div class="info-container">
        <span class="like-text"> Comment </span>
        <button class="like-button"></button>
        <span class="like-text"> Put a Like  </span>
        <button @click="toggleLike(photo.photo_Id)"
                class="like-button"
                :class="photo.liked ? 'text-primary': '' "    >
        </button>
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
  font-size: 25px;            /* Dimensione del testo più piccola */
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
  font-size: 30px; /* Dimensione dell'icona */

}

.text-primary {
  background-color: #1877f2 !important; /* Colore blu per il like attivo, simile a Facebook */
}

.like-button:hover i {
  color: #6c757d;  /* Cambia colore al passaggio del mouse */
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
  display: flex;
  position: relative; /* Posizione relativa per il posizionamento degli elementi */
}

.fixed-bottom {
  position: absolute;
  bottom: 0; /* Posiziona il div in basso */
  left: 0; /* Inizia dall'angolo sinistro */
  width: 100%; /* Occupa tutta la larghezza dello schermo */
  justify-content: flex-start; /* Allinea gli elementi a sinistra */
  padding: 10px; /* Spaziatura interna */
  margin: 0;
}
.comments-container{
  width: 100%;
}


</style>