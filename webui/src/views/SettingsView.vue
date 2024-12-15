<script>
import defaultProfileImage from '../assets/images/default-profile.jpg';
import GoToButton from '../components/GoToButton.vue';

export default {
    data() {
        return {
            Username: sessionStorage.getItem("Username"),
            Identifier: sessionStorage.getItem("Auth"),
            PhotoCode: sessionStorage.getItem("PhotoProfile"),
            PhotoURL: '',
            ErrorMessage: '',
            isEditingName: false, // stato per sapere se stiamo modificando il nome
            newName: this.Username,
        }
    },
    methods: {
        async GetPhotoProfile(){
            try {
                let response = await this.$axios.get(`/images/${this.PhotoCode}/photo`, {
                    headers: {
                        Authorization: `Bearer ${this.Identifier}`, 
                    },
                    responseType: 'blob',
                });
                
                this.PhotoURL = URL.createObjectURL(response.data);

            } catch(error) {
                if (error.response) {
                    this.ErrorMessage = `Error: ${error.response.status} - ${error.response.data}`;
                } else if (error.request) {
                    // Se la richiesta è stata fatta ma non c'è risposta
                    this.ErrorMessage = "Network error, please try again.";
                } else {
                    // Qualcosa è andato storto nella configurazione della richiesta
                    this.ErrorMessage = "Unexpected error: " + error.message;
                }
            }
        },

        onImageClick() {
            // simula il click sull'input nascosto
            this.$refs.fileInput.click();
        },

        handleFileChange(event) {
            const file = event.target.files[0];
            if (file) {
                this.SetPhotoProfile(file);
            }
        },

        async SetPhotoProfile(file){
            try {

                if (file.type !== 'image/png') {
                    this.ErrorMessage = "Only PNG images are allowed.";
                    return;
                }

                let response = await this.$axios.put(`/users/${this.Identifier}/photoProfile`, file, {
                    headers: {
                        Authorization: `Bearer ${this.Identifier}`,
                        "Content-Type": "image/png",
                    },
                });

                this.PhotoURL = URL.createObjectURL(file);
                sessionStorage.setItem("PhotoProfile", response.data.imageCode);
                this.PhotoCode = response.data.imageCode;

            } catch(error) {

                if (error.response) {
                    this.ErrorMessage = `Error: ${error.response.status} - ${error.response.data}`;
                } else if (error.request) {
                    // Se la richiesta è stata fatta ma non c'è risposta
                    this.ErrorMessage = "Network error, please try again.";
                } else {
                    // Qualcosa è andato storto nella configurazione della richiesta
                    this.ErrorMessage = "Unexpected error: " + error.message;
                }
            }
        },
        
        onNameClick(){
            this.isEditingName = true;
        },

        handleNameChange(event) {
            const newName = event.target.value;
            this.newName = newName;

            this.updateName(newName);
        },

        async updateName(newName) {
            try {
                let response = await this.$axios.put(`/users/${this.Identifier}/username`, { username: newName }, {
                    headers: {
                        Authorization: `Bearer ${this.Identifier}`,
                        "Content-Type": "application/json",
                    },
                });
                this.Username = newName;
                sessionStorage.setItem("Username", newName);
                this.isEditingName = false; // esci dalla modalità di editing
            } catch (error) {
                this.ErrorMessage = error.response ? `Error: ${error.response.status} - ${error.response.data}` : "Unexpected error: " + error.message;
            }
        },

    },

    mounted() {
        if (this.PhotoCode) {
            this.GetPhotoProfile();
        } else {
            this.PhotoURL = defaultProfileImage;

        }
    },

}
</script>

<template>
    <center>
        <div class="container-profile">
            <div class="image" @click="onImageClick">
                <img :src="PhotoURL" alt="Profile Picture" class="profile-image" />
                <span class="edit">Edit</span>
            </div>
            <input type="file" ref="fileInput" style="display: none" @change="handleFileChange" accept="image/png"/>
            <div class="name" @click="onNameClick">
                <h1 v-if="!isEditingName">{{ Username }}</h1>
                <input v-else type="text" class="nameInput" v-model="newName" @blur="handleNameChange" /> 
                <span v-if="!isEditingName" class="edit">Edit</span>
            </div>
            
        </div>
        <ErrorBox v-if="ErrorMessage" :msg="ErrorMessage"></ErrorBox>
    </center>
</template>


<style scoped>

.container-profile {
    margin-top: 100px;
}

.image {
    position: relative; /* necessario per posizionare il testo "Edit" sopra */
    border-radius: 100%;
    width: 400px;
    height: 400px;
    overflow: hidden; 
    -webkit-box-shadow: 0px 0px 5px 4px rgba(0,0,0,0.75);
    -moz-box-shadow: 0px 0px 5px 4px rgba(0,0,0,0.75);
    box-shadow: 0px 0px 5px 4px rgba(0,0,0,0.75);
    cursor: pointer;
}

.image img {
    width: 100%;
    height: 100%; 
    object-fit: cover;
    pointer-events: none;
    transition: opacity 0.3s ease;
}

.image:hover img {
    opacity: 0.5;
}

.image .edit {
    position: absolute;
    top: 50%;
    left: 50%; 
    transform: translate(-50%, -50%);
    color: white;
    font-size: 30px; 
    font-family: Arial, sans-serif;
    font-weight: bold;
    text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.7);
    opacity: 0; /* inizia invisibile */
    transition: opacity 0.3s ease;
    pointer-events: none; /* evita conflitti con il click */
}

.image:hover .edit {
    opacity: 1; /* rende visibile quando ci vai sopra */
}


.name {
    margin-top: 20px;
}

.name {
    margin-top: 20px;
    position: relative; /* necessario per posizionare la scritta "Edit" */
    display: inline-block; 
}

.name h1 {
    font-size: 70px;
    font-family: sans-serif;
    transition: text-decoration 0.3s ease; 
    cursor: pointer; 
}
.name h1:hover {
    opacity: 0.5;
}

.name .edit {
    position: absolute;
    top: 0;
    left: 100%;
    font-size: 16px; 
    color: gray;
    font-family: Arial, sans-serif;
    font-weight: normal;
    opacity: 0; /* inizia invisibile */
    transform: translateY(10px); /* effetto spostamento iniziale */
    transition: opacity 0.3s ease, transform 0.3s ease; 
}

.name:hover .edit {
    opacity: 1; /* mostra il testo "Edit" */
    transform: translateY(0); /* eiporta la scritta alla posizione originale */
}
</style>