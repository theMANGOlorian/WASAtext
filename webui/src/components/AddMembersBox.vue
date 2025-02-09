<template>
    <div class="user-list-overlay" v-on:click.self="close">
        <div class="user-list">
            <h2>Aggiungi utenti al gruppo</h2>
            <ul>
                <li 
                    v-for="user in usersList" 
                    :key="user.username"  
                    v-on:click="addToGroup(user.username)">
                    <div class="contact-item">
                        <span class="contact-name">{{ user.username }}</span>
                    </div>
                </li>
            </ul>
            <button v-on:click="close" class="closeButton">Close</button>
        </div>
    </div>
    
</template>

<script>
export default {
    props: {
        groupId: {
            type: String,
            required: true
        }
    },

    data() {
        return {
            auth: '',
            usersList: [],
        };
    },

    beforeMount() {
        this.auth = sessionStorage.getItem('Auth');
        this.fetchUsersList();
    },

    methods: {
        async fetchUsersList() {
            try {
                const response = await this.$axios.get(`/users`, {
                    headers: { Authorization: `Bearer ${this.auth}` },
                });

                this.usersList = response.data.users.map(username => ({
                    username
                })) || [];
            } catch (error) {
                console.error("Errore nel recupero degli utenti:", error.message);
                alert("Error Message: " + error.message);
            }
        },

        async addToGroup(username) {
            try {
                await this.$axios.put(`/conversations/${this.groupId}/member`, 
                    { username },
                    {
                        headers: { Authorization: `Bearer ${this.auth}` },
                    }
                );
                alert("L'utene " + username + " è stato aggiunto con successo");
            } catch (error) {
                console.error(`Errore nell'aggiunta di ${username} al gruppo:`, error.message);
                alert("Error Message: " + error.message);
            }
        },

        close() {
            this.$emit('close'); // Chiude il popup
        }
    }
};
</script>


<style scoped>
.user-list-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.user-list {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
    width: 350px;
    padding: 20px;
    text-align: center;
}

.user-list h2 {
    margin-bottom: 20px;
    font-size: 1.5rem;
    color: #333;
}

.user-list ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.user-list li {
    margin: 10px 0;
    background: #f9f9f9;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.3s, transform 0.2s;
    padding: 10px;
}

.user-list li:hover {
    background: #65ee61;
    transform: scale(1.02);
}

.contact-item {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    font-weight: 500;
    color: #333;
}
.closeButton {
    margin-top: 10px;
    padding: 5px 30px;
    font-size: 1rem;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: 0.1s ease ;
}

.closeButton:hover {
    background-color: #b30000;
}
</style>
