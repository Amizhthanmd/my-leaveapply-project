<template>
  <header class="header">
		<h1 class="logo">###</h1>
      <ul class="main-nav">
          <li><a class="link" @click="home" href="#/"> EMP LEAVE BALANCE INFO <i class="bi bi-info-circle"></i>  </a> </li>
          <li> <a class="link"  @click="leavestatus" href="#/Submitted_Details" > LEAVE STATUS <i class="bi bi-clock-history"></i> </a> </li>
          <li><a class="link" @click="approve" href="#/HR_Approve" > APPROVAL PORTAL <i class="bi bi-check-circle"></i> </a> </li>
      </ul>
	</header> 
  <component :is="currentView" />
<div>
  <loading 
    :active='isLoading' 
    :is-full-page="fullPage" 
    :loader='loader'
    :color='color' />
</div>
</template>

<script>

  import Employee_Info from './Employee_Info.vue'
  import Submitted_Details from './Submitted_Details.vue'
  import HR_Approve from './HR_Approve.vue'
  import Main_Grid from './Main_Grid.vue'
  import Grid_View from './Grid_View.vue'
  import Form_View from './Form_View.vue'
  import Loading from 'vue-loading-overlay';
  import 'vue-loading-overlay/dist/css/index.css';
  import '@progress/kendo-ui';
  import '@progress/kendo-theme-default/dist/all.css'



  const routes = {
    '/': Employee_Info,
    '/Submitted_Details': Submitted_Details,  
    '/HR_Approve': HR_Approve,
    '/Main_Grid': Main_Grid, 
    '/Grid_View': Grid_View,  
    '/Form_View': Form_View  

  }
  
  export default {
    components:{
      Loading
    },
    data() {
      return {
        currentPath: window.location.hash,
        isLoading: false,
        fullPage: true,
        loader: 'bars',
        color: 'rgb(17, 107, 139)', 
      }
    },
    computed: {
      currentView() {
        return routes[this.currentPath.slice(1) || '/']
      }
    },
    methods: {
      home(){
        this.isLoading = true;
        setTimeout(() => {this.isLoading = false}, 800)
      },
      leavestatus(){
        this.isLoading = true;
        setTimeout(() => {this.isLoading = false}, 800)
      },
      approve(){
        this.isLoading = true;
        setTimeout(() => {this.isLoading = false}, 800)
      },
    },
    mounted() {
      window.addEventListener('hashchange', () => {
        this.currentPath = window.location.hash
      })
    }
  }
  
</script>

<style>
   
* {
	box-sizing: border-box;
}
body{
	font-family: 'Montserrat', sans-serif;
	line-height: 1.6;
	margin: 0;
	min-height: 100vh;
}
ul{
  margin: 0;
  padding: 0;
  list-style: none;
}
.logo{
	color: #ffffff;
}
.logo {
	margin-left: 20px;
	font-size: 1.45em;
}
.main-nav {
	margin-top: 5px;
}
.link{
  display: inline-block;
  position: relative;
  color: #fdfdfd;
  text-decoration: none;
	font-size:16px;
  font-weight: bold;
}
.link:after {
  content: '';
  position: absolute;
  width: 100%;
  transform: scaleX(0);
  height: 2px;
  bottom: 0;
  left: 0;
  background-color: #0a7eb8;
  transform-origin: bottom right;
  transition: transform 0.25s ease-out;
}
.link:hover:after {
  transform: scaleX(1);
  transform-origin: bottom left;
}
.main-nav a {
	padding: 10px 15px;
	text-transform: uppercase;
	text-align: center;
	display: block;
}
.main-nav a:hover {
	color: #ffffff;
}
.header {
	padding-top: .5em;
	padding-bottom: .5em;
	border: 1px solid #a2a2a2;
	background-color: #11abca;
	-webkit-box-shadow: 0px 0px 14px 0px rgba(101, 101, 101, 0.75);
	-moz-box-shadow: 0px 0px 14px 0px rgba(93, 93, 93, 0.75);
	box-shadow: 0px 0px 14px 0px rgba(94, 94, 94, 0.75);
}
@media (min-width: 769px) {
	.header,
	.main-nav {
		display: flex;
	}
.header {
		flex-direction: column;
		align-items: center;
    width: 100%;
		margin: 0 auto;
}
}
@media (min-width: 1025px) {
	.header {
		flex-direction: row;
		justify-content: space-between;
	}
}
</style>
