<template>
    <br>
    <div class="headerlink">
        <a href="#/Form_View" > Add New Announcement </a>
    </div>
    <div class="k-grid-header">
        <kendo-grid :data-source="receive" ref="maingrid"
                    :pageable="pageable"
                    :filterable="true"
                    :noRecords="true"
                    :height="300">
        <kendo-grid-column :field="'Title'"  
                           :title="'TITLE'"></kendo-grid-column>
        <kendo-grid-column :field="'Audience'" 
                           :title="'AUDIENCE'" ></kendo-grid-column>
        <kendo-grid-column :field="'Types'" 
                           :title="'TYPES'" ></kendo-grid-column>
        <kendo-grid-column :field="'Postedon'" 
                           :title="'POSTED ON'" ></kendo-grid-column>
        <kendo-grid-column :field="'Postedby'" 
                           :title="'POSTED BY '" ></kendo-grid-column>
        <kendo-grid-column :command="[{name:'Download',click: view}]"
                           :title="'ATTACHMENTS'" ></kendo-grid-column>
        </kendo-grid>
  </div>
</template>

<script>
        import { Grid, GridColumn } from '@progress/kendo-grid-vue-wrapper';
        import axios from 'axios';
        import $ from 'jquery';

    export default {
    components: {  
        'kendo-grid'          : Grid,
        'kendo-grid-column'   : GridColumn,
        // 'kendo-dropdownlist'  : DropDownList,
        },
        data: function(){
            return{
                receive : {
                    Title : "",
                    Audience : "",
                    Types : "",
                    Postedon : "",
                    Postedby : "",
                    Attachment : "",
                    },
                    pageable: { pageSizes: [1, 5, 10, 20],
                                buttonCount: 5,
                                pageSize: 10,
                            },
            }
        },
        created(){
            this.get()
        },
        methods: {
            get(){
                axios.get('http://localhost:8007/getgridview',{})
                .then((response) => {
                    this.receive = response.data
                    this.get
                })
                .catch((error) => {
                    console.log(error)
                })
            },
            view(e){
              var gridWidget = this.$refs.maingrid.kendoWidget();
              var tr = $(e.target).closest("tr");
              gridWidget.dataItem(tr);
              console.log("Attachment : ",gridWidget.dataItem(tr).Attachment)
              var file = gridWidget.dataItem(tr).Attachment
                if(file ==""){
                    alert("No Attachment found ")
                    } else {
                    var link = './images/'+file
                    axios({
                    url: link,
                    method: 'GET',
                    responseType: 'blob',
                    }).then((response) => {
                    var fileURL = window.URL.createObjectURL(new Blob([response.data]))
                    var fileLink = document.createElement('a')
                    fileLink.href = fileURL
                    fileLink.setAttribute('download', file)
                    document.body.appendChild(fileLink)
                    fileLink.click()
                    })
                }
              console.log("File Name : ",file)
            },
        },
    }

</script>

<style>
.headerlink{
    text-align: center;
}
</style>


