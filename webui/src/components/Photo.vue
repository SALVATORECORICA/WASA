<script>
export default {

  props: {
    photos: {},
  },

  data() {
    return {
      errormsg: null,
      commentModal: false,
      openedPhoto: null,
      newComment: "",
      requester: Number(localStorage.getItem('token'))
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
            await this.$axios.put("/users/" + this.requester + "/photos/" + p.photo_Id + "/likes/" + localStorage.getItem('token'), {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              },
            });
            p.nLikes = p.nLikes + 1
          } catch (e) {
            p.liked = !p.liked;
            this.errormsg = e.toString();
          }
        } else {
          try {
            await this.$axios.delete("/users/" + this.requester + "/photos/" + p.photo_Id + "/likes/" + localStorage.getItem('token'), {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              },
            });
            p.nLikes = p.nLikes - 1
          } catch (e) {
            p.liked = !p.liked;
            this.errormsg = e.toString();
          }
        }
      }
    },
    async openComments(photo) {
      // Recupera nuovamente i dati della foto per aggiornare i commenti
      try {
        const response = await this.$axios.get("/users/" + this.requester + "/photos/" + photo.photo_Id, {
          headers: {
            "Authorization": `Bearer ${localStorage.getItem('token')}`,
          },
        });

        // Aggiorna la foto selezionata con i nuovi dati
        this.openedPhoto = response.data;
        this.commentModal = true;
      } catch (e) {
        console.error(e); // Gestione dell'errore nella richiesta GET
      }
    },

    closeComments() {
      this.openedPhoto = null;
      this.commentModal = false;
      this.newComment = "";
    },
    async addComment() {
      try {
        // Invia il commento con una richiesta POST
        await this.$axios.post(`/users/${this.requester}/photos/${this.openedPhoto.photo_Id}/comments`, {
          comment: this.newComment,
        }, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        });

        await this.fetchPhoto();

          // Reset del campo di input
          this.newComment = "";
        } catch (e) {
          console.error(e); // Gestione dell'errore nella richiesta GET
        }
    },

    async deleteComment(comment) {
      try {
        // Invia il commento con una richiesta DELETE
        await this.$axios.delete(`/users/${this.requester}/photos/${this.openedPhoto.photo_Id}/comments/${comment.comment_id}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        });

        // Recupera nuovamente i dati della foto per aggiornare i commenti
        await this.fetchPhoto();
      } catch (e) {
        console.error(e); // Gestione dell'errore nella richiesta DELETE
      }
    },

    

    async fetchPhoto() {
      try {
        const response = await this.$axios.get(`/users/${this.requester}/photos/${this.openedPhoto.photo_Id}`, {
          headers: {
            "Authorization": `Bearer ${localStorage.getItem('token')}`,
          },
        });

        // Aggiorna la foto selezionata con i nuovi dati
        this.openedPhoto = response.data;
      } catch (e) {
        console.error(e); // Gestione dell'errore nella richiesta GET
      }
    },
    async deletePhoto(photo_Id) {
      try {
        await this.$axios.delete("/users/" + this.requester + "/photos/" + photo_Id, {
          headers: {
            "Authorization": `Bearer ${localStorage.getItem('token')}`
          },
        });
        this.$emit("photoDeleted");
      } catch (e) {
        console.log(e.toString());
      }
    }
  },

}

</script>

<template>
  <div v-if="commentModal" @click.self="closeComments" class="overlay-background">
    <div class="overlay">
      <div class="comments-container">
        <div v-for="comment in openedPhoto.comments" :key="comment.comment_id">
          <span> {{ comment.user.nickname}}</span>
          <span> {{comment.comment}}</span>
          <button  v-if=" (requester === openedPhoto.owner.id) ||
                           (requester === comment.user.id)"
                   class="button-delete"
                   @click="deleteComment(comment)">
                   Delete Comment
          </button>
        </div>
      </div>
      <div class="fixed-bottom">
        <input placeholder="Insert your comment" v-model="newComment">
        <button @click="addComment"> Add comment</button>
      </div>
    </div>
  </div>

  <div v-if="!commentModal">
    <div v-for="photo in photos" :key= " photo.photo_Id" class="photo">
      <h2>{{ photo.owner.nickname }} </h2>
      <img :src= "'data:image/png;base64, ' + photo.image" />
      <div class="info-container">
        <span class="like-text"> Comment </span>
        <button class="like-button" @click="openComments(photo)"></button>
        <span class="like-text"> Put a Like  </span>
        <button @click="toggleLike(photo.photo_Id)"
                class="like-button"
                :class="photo.liked ? 'text-primary': '' "    >
        </button>
        <span class="like-text"> Likes:  {{ photo.nLikes }} </span>
        <span  v-if=" (requester === photo.owner.id) "
               class="like-text delete-photo"
               @click="deletePhoto(photo.photo_Id)">
          Delete
        </span>
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
  border-radius: 10px;
  flex-direction: column;
  display: flex;
  position: relative; /* Posizione relativa per il posizionamento degli elementi */
}
.comments-container {
  flex: 1;                    /* Occupa lo spazio disponibile */
  overflow-y: auto;           /* Rende la lista dei commenti scrollabile */
  margin-bottom: 50px;        /* Lascia spazio per l'input */
}

.fixed-bottom {
  position: absolute;
  bottom: 0; /* Posiziona il div in basso */
  left: 0; /* Inizia dall'angolo sinistro */
  width: 100%; /* Occupa tutta la larghezza dello schermo */
  justify-content: space-between;
  padding: 10px; /* Spaziatura interna */
  margin: 0;
  display: flex;
}
.comments-container div {
  padding: 10px;
  border-bottom: 1px solid #eaeaea;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  position: relative; /* Posizionamento relativo per il pulsante */
}



/* Stile per il nickname dell'utente in alto */
.comments-container span:first-child {
  font-weight: bold;
  color: #333; /* Colore più scuro per il nome */
  margin-bottom: 5px;
}

/* Stile per il testo del commento in basso */
.comments-container span:nth-child(2) {
  color: #555;
  margin-bottom: 10px; /* Spazio sotto il commento */
}

.delete-photo {
  color: deeppink;
  cursor: pointer;
  transition: transform 0.2s ease, color 0.2s ease;
}

.delete-photo:hover {
  transform: translateY(-3px); /* Alza leggermente il testo */
  color: red; /* Cambia il colore al passaggio del mouse */
}



</style>