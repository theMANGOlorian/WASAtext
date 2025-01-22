<script>
export default {
    data() {
        return {
            selectedConversation: null, // La conversazione attualmente selezionata
        };
    },

    methods: {
        handleConversationSelected(conversation) {
            this.selectedConversation = conversation;
        },

        nullSelectedConversation() {
            this.selectedConversation = null;
        },
    },
    
};
</script>

<template>
    <div id="app" class="webchat-container">
      <!-- Sidebar con la lista delle chat -->
      <aside class="chat-list">
        <ConversationBox @conversation-selected="handleConversationSelected"></ConversationBox>
      </aside>
  
      <!-- Sezione principale per la chat selezionata -->
      <main class="chat-content">
        <ChatBox :conversation="selectedConversation" v-on:rmSelectedConversation="nullSelectedConversation"></ChatBox>
      </main>
    </div>
</template>

<style scoped>

.webchat-container {
    display: flex;
    height: 100vh;
    overflow: hidden;
}

.chat-list {
    width: 25%;
    background-color: #f5f5f5;
    border-right: 1px solid #ddd;
    overflow-y: auto;
    
}

.chat-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    background-color: #fff;
}
</style>
  