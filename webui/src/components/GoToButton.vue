<template>
    <button @click="goToPage" class="go-to-button">
      {{ text }}
    </button>
</template>
  
<script>
export default {
  props: {
    text: {
      type: String,
      required: true
    },
    path: {
      type: String,
      required: true
      },
  },
  methods: {
    goToPage() {
      let p = ''
      if (this.path == 'chats'){
        p = '/users/'+sessionStorage.getItem('Auth')+'/chats';
      } else if (this.path == 'settings') {
        p = '/users/'+sessionStorage.getItem('Auth')+'/settings';
      }
      this.$router.push(p);
    }
  },
};
</script>
  
  <style scoped>
  .go-to-button {
    background-color: rgb(37, 145, 4);
    color: rgb(255, 255, 255);
    font-size: 16px;
    padding: 15px 60px;
    border: none;
    cursor: pointer;
    transition: background-color 0.3s;
    font-family: Verdana, Geneva, Tahoma, sans-serif;
  }
  .go-to-button:hover {
    background-color: #70b615;
  }
  </style>
  