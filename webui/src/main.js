import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorBox from './components/ErrorBox.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import GoToButton from './components/GoToButton.vue';
import conversationBox from './components/ConversationBox.vue';
import ChatBox from './components/ChatBox.vue';
import ContextMenu from './components/ContextMenu.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App);
app.config.globalProperties.$axios = axios;
app.component("ErrorBox", ErrorBox);
app.component("LoadingSpinner", LoadingSpinner);
app.component("GoToButton", GoToButton);
app.component("ConversationBox", conversationBox);
app.component("ChatBox",ChatBox);
app.component("ContextMenu", ContextMenu);
app.use(router);
app.mount('#app');
