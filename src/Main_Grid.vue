<template>
        <br>
        <div class="k-grid-header">
            <kendo-grid :data-source="receive"
                        :pageable="pageable"
                        :filterable="true"
                        :noRecords="true"
                        :height="300">
            <kendo-grid-column :field="'Name'"  
                               :title="'NAME'"></kendo-grid-column>
            <kendo-grid-column :field="'Gender'" 
                               :title="'GENDER'" ></kendo-grid-column>
            <kendo-grid-column :field="'Skills'" 
                               :title="'SKILLS'" ></kendo-grid-column>
            <kendo-grid-column :field="'Emailid'" 
                               :title="'EMAIL ID'" ></kendo-grid-column>
            </kendo-grid>
      </div>
</template>
    
<script>
    import { Grid, GridColumn } from '@progress/kendo-grid-vue-wrapper';
    import axios from 'axios';

    export default {
            components: {  
            'kendo-grid'          : Grid,
            'kendo-grid-column'   : GridColumn,
         // 'kendo-dropdownlist'  : DropDownList,
    },
    data: function(){
        return{
            receive: {
                Name : '',
                Gender : '',
                Skills : [],
                Emailid : '',
            },
            pageable: { pageSizes: [1, 5, 10, 20],
                buttonCount: 5,
                pageSize: 10,
              },
        }
    },
    created() {
          this.get()
    },
    methods: {
        get(){
            axios.get("http://localhost:8007/requestempdetails", {})
              .then((response) => {
              try{
              for(var i=0; i<response.data.length; i++){
              if(response.data[i].Skills){
                var page=response.data[i].Skills
                var val=""
                page.forEach((x)=>{
                val+=x+", "
                })
                response.data[i]["Skills"]=val
                }
              }
              }catch(error){
                console.log("Error",error)
              }
              console.log(response.data)
              this.receive =response.data
            })
            },
    }
    
}

</script>
    
<style>
    
</style>



