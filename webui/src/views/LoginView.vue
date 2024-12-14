<template>
    <center>
        <div class="title-box">
            <h1 class="title" style="margin-right: 100px;"><span style="color: purple">Talk</span> with <span style="color: red">friends</span></h1>
            <h1 class="title" style="margin-left: 200px;"> on <span style="color: rgb(37, 145, 4)">WasaText</span>!</h1>
        </div>
        <div class="input-username-box">
            <input class="input-text" type="text" placeholder="Enter your username" v-model="Username" />
        </div>
        <div class="button-box">
            <button id="submit"class="submit" @click="LoginButtonPressed" :disabled="!validateUsername()" :style="{backgroundColor: validateUsername() ? 'rgb(37, 145, 4)' : 'rgba(37, 145, 4,0.5)'}">Login</button>
        </div>
        <ErrorBox v-if="ErrorMessage" :msg="ErrorMessage"></ErrorBox>
    </center>

</template>

<script>


export default {
    data() {
        return {
            Username: '',
            ErrorMessage: '',
            isUsernameValid: false,
        }
    },

    methods: {    
        validateUsername() {
            const regex = /^[0-9a-zA-Z]{3,25}$/;            
            // Verifica che il nome utente sia valido
            if (!regex.test(this.Username)) {
                this.ErrorMessage = "Username must contain only [0-9a-zA-Z] and be between 3 to 25 characters.";
                return false;
            }

            this.ErrorMessage = ''; // Rimuovi il messaggio di errore
            return true;
        },


        async LoginButtonPressed(){
            try {
                let doLogin_response = await this.$axios.post("/session",{Username:this.Username}); // Username è il campo nella request del json, this.Usename fa riferimento al v-model
                sessionStorage.setItem("Auth", doLogin_response.data.identifier);
                sessionStorage.setItem("Username",doLogin_response.data.username);
                sessionStorage.setItem("PhotoProfile", doLogin_response.data.photoCode)

                this.$router.push("/users/" + sessionStorage.getItem("Auth") + "/chats");

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
        }

    }
};
</script>

<style scoped>

/* Stile del titolo */
.title {
    font-family: sans-serif;
    font-size: 40px;
    font-weight: bold;
    color: black; /* Per rendere il titolo visibile sullo sfondo */
}
.title-box {
    margin-top: 200px;
}

.input-username-box{
    margin-top: 50px;

}

.input-text:hover{
    width: 400px;
}

.input-text {
    width: 250px;
    text-align: center;
    font-size: 20px;
    border: none;
    border-bottom: solid 2px rgb(79, 10, 107);
    outline: none;
    transition: 0.9s ease;
}

.button-box {
    margin-top: 30px;
}
.button-box .submit {
    width: 200px;
    border-radius: 20px;
    font-size: 25px;
    color: white;
    background-color: rgb(3, 77, 13);
    transition: 0.5s ease;
}

.button-box .submit:hover{
    background-color: white;
    color: rgb(0, 0, 0);
    border-color: rgb(0, 0, 0);
}





/**  Spec for smartphone screen size */
@media (max-width: 767.98px) { 

}


</style>
