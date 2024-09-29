<script>
import errorMsg from "../components/ErrorMsg.vue";

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			isPhotoSelected: false,
      photo: null,
      inputKey: Date.now(),
      photos: [],
		}
	},

  props: ["nickname", "id"],

  methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
    async deselectPhoto(){
      this.inputKey = Date.now();
      try {
        await  this.$nextTick();
      } catch (e){
        this.errormsg = e.toString();
      }
      this.photo = null;
      this.isPhotoSelected = false
    },
    selectPhoto(event){
      const file = event.target.files[0]; // Ottieni il primo file selezionato
      if (file) {
        this.photo = file; // Memorizza il file se esiste
        this.isPhotoSelected = true


      } else {
        this.photo = null; // Resetta se non c'è file
      }
    },
    async uploadPhoto() {
      let fileInput = document.getElementById('fileUploader')

      const file = fileInput.files[0];
      const reader = new FileReader();

      reader.readAsArrayBuffer(file);

      reader.onload = async () => {
        // Post photo: /users/:id/photos
        let response = await this.$axios.post("/users/" + this.$route.params.id + "/photos", reader.result, {
          headers: {
            'Content-Type': file.type
          },
        })
      }
    },
    async stream(){
        try {
          this.errormsg = null
          let response = await this.$axios.get("/users/" + this.id+ "/home",
              {
                headers:{
                  'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
              })
          this.photos = response.data

        } catch (e) {
          this.errormsg = e.toString()
        }
      }
    },
	mounted() {
		//this.refresh()
    this.stream()

	}
}
</script>

<template>
	<div class="flex-grow-1">
    <nav class="navbar navbar-expand navbar-light bg-light sticky-top" style="background-color: transparent; border: none; box-shadow: none;">
      <div class="container-fluid">
        <span class="navbar-text h2">Welcome {{ this.nickname }}</span>
        <div class="ms-auto">
          <div class="btn-group me-2">
            <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
              Refresh
            </button>
          </div>
          <div class="btn-group me-2">
            <!--Uso di key per deselezionare il file, aggiornando la key il framework forza un nuovo rendering -->
            <input type="file" ref="fileInput" class="btn btn-sm btn-outline-secondary"  @input="selectPhoto" accept=".jpg, .png" :key="inputKey">
            <button  @click="deselectPhoto">Deselect Photo</button>
            <button  v-if="isPhotoSelected" @click="uploadPhoto">Upload selected Photo</button>
          </div>
          <div class="btn-group me-2">
            <button type="button" class="btn btn-sm btn-outline-primary">
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
