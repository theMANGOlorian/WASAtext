<template>
    <div 
        class="context-menu"
        v-if="visible" 
        :style="{ top: position.y + 'px', left: position.x + 'px' }"
        @click.stop
    >
        <ul>
            <li @click="onOptionClick(option)" v-for="option in options" :key="option">
                {{ option }}
            </li>
        </ul>
    </div>
</template>

<script>
export default {
    data() {
        return {
            visible: false,
            position: { x: 0, y: 0 },
            options: [],
        };
    },
    methods: {
        openMenu(x, y, options) {
            this.position = { x, y };
            this.options = options;
            this.visible = true;
        },
        closeMenu() {
            this.visible = false;
        },
        onOptionClick(option) {
            this.$emit("option-click", option);
            this.closeMenu();
        },
        getMenuPosition(){
            return {
                x: this.position.x,
                y: this.position.y,
            };
        },
    },
};
</script>

<style scoped>
.context-menu {
    position: absolute;
    background-color: rgba(255, 255, 255, 0.911);
    border: 1px solid #ccc;
    border-radius: 10px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    width: fit-content;
}

.context-menu ul {
    list-style: none;
    margin: 0;
    padding: 0;
}

.context-menu li {
    padding: 10px;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.context-menu li:hover {
    background-color: #d6d3d3;
}
</style>
