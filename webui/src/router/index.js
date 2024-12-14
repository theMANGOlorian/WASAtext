import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'
import SettingsView from '../views/SettingsView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', redirect: '/session'},
		{path: '/session', component: LoginView},
		{path: '/users/:Identifier/home', component: HomeView},
		{path: '/users/:Identifier/chats', component: ChatView},
		{path: '/users/:Identifier/settings', component: SettingsView},
	]
})

export default router
