<script>
import errorMsg from "../components/ErrorMsg.vue";
import Photo from "../components/Photo.vue";

export default {
  components: {Photo},
  data: function () {
    return {
      errormsg: null,
      loading: false,
      isPhotoSelected: false,
      photo: null,
      inputKey: Date.now(),
      photos: [],
      changedName: "",
      localnickname : localStorage.getItem('nickname'),
      id : Number(localStorage.getItem('token'))
    }
  },

  props: ["nickname",],

  mounted() {
    //this.refresh()
    this.stream();

  },

  methods: {
    async refresh() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.get("/");
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async deselectPhoto() {
      this.inputKey = Date.now();
      try {
        await this.$nextTick();
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.photo = null;
      this.isPhotoSelected = false
    },
    selectPhoto(event) {
      const file = event.target.files[0]; // Ottieni il primo file selezionato
      if (file) {
        this.photo = file; // Memorizza il file se esiste
        this.isPhotoSelected = true


      } else {
        this.photo = null; // Resetta se non c'è file
      }
    },
    async uploadPhoto() {
      let fileInput = this.$refs.fileInput;

      const file = fileInput.files[0];
      const reader = new FileReader();

      reader.readAsArrayBuffer(file);

      reader.onload = async () => {
        try {
          // Post photo: /users/:id/photos
          await this.$axios.post("/users/" + localStorage.getItem('token') + "/photos", reader.result, {
            headers: {
              'Content-Type': file.type,
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
          });
          // Deseleziona il file dopo l'upload
          await this.deselectPhoto();
        } catch (e) {
          // Gestione degli errori
          this.errormsg = e.toString();
        }
      };
    },


    async stream() {
      try {
        this.errormsg = null
        let response = await this.$axios.get("/users/" + this.id + "/home",
            {
              headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
              }
            })
        this.photos = response.data

      } catch (e) {
        this.errormsg = e.toString()
      }
    },
    logout() {
      this.$emit("logout")

    },
    async changeName() {
      if (this.changedName.length >= 3 && this.changedName.length <= 16) {
        try {
          // Effettua la richiesta PUT
          console.log(this.id)
          await this.$axios.put("/users/" + this.id, {
            "nickname": this.changedName
          }, {
            headers: {
              "Authorization": `Bearer ${localStorage.getItem('token')}`
            }
          });

          // Aggiorna il nickname e resetta il campo
          localStorage.setItem("nickname", this.changedName)
          this.changedName = "";
          this.$emit("changeName", this.localnickname);
        } catch (e) {
          console.log("Errore durante il cambio del nickname:", e);  // gestisci eventuali errori
        }
      } else {
        console.log("Il nickname deve essere lungo tra 3 e 16 caratteri");
      }
    }

  }
}
</script>

<template>
	<div class="flex-grow-1 ">
    <nav class="navbar navbar-expand navbar-light bg-light" style="background-color: transparent; border: none; box-shadow: none;">
      <div class="container-fluid ">
        <span class="navbar-text h2">Welcome {{ localnickname }} </span>
        <div class="ms-auto">
          <div class="btn-group me-2">
            <!--Uso di key per deselezionare il file, aggiornando la key il framework forza un nuovo rendering -->
            <input type="file" ref="fileInput" class="btn btn-sm btn-outline-secondary"  @input="selectPhoto" accept=".jpg, .png" :key="inputKey">
            <button  @click="deselectPhoto">Deselect Photo</button>
            <button  v-if="isPhotoSelected" @click="uploadPhoto" class="upload">Upload selected Photo</button>
          </div>
          <div class="btn-group me-2">
            <input v-model="changedName" placeholder="Insert your new Nickname">
            <button  @click="changeName" type="button" class="btn btn-sm btn-outline-primary">
              Change Name
            </button>
            <button @click="logout" type="button" class="btn btn-sm btn-outline-primary">
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
  <Photo :photos ="photos"  :id="id" ></Photo>
</template>

<style>
.upload{
  background-color: red;
}
</style>
