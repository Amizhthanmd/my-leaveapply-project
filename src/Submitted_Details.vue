<template>
<br>
<div class="head">
    <h4> LEAVE APPLIED DETAILS AND STATUS </h4>
</div>
<hr>
<div id="kgridheader">
      <kendo-grid :data-source="receive" 
                  :pageable="pageable"
                  :filterable="true"
                  :noRecords="true"
                  :height="300">
      <kendo-grid-column :field="'EmployeeId'"
                         :title="'EMPLOYEE ID'" ></kendo-grid-column>
      <kendo-grid-column :field="'Gender'"
                         :title="'GENDER'" ></kendo-grid-column>
      <kendo-grid-column :field="'LeaveType'"
                         :title="'LEAVE TYPE'" ></kendo-grid-column>
      <kendo-grid-column :field="'Leavebalance'"
                         :title="'LEAVE BALANCE'" ></kendo-grid-column>
      <kendo-grid-column :field="'FromDate'"
                         :title="'FROM DATE'" ></kendo-grid-column>
      <kendo-grid-column :field="'ToDate'"
                         :title="'TO DATE'" ></kendo-grid-column>
      <kendo-grid-column :field="'NoofDays'"
                         :title="'NO OF DAYS'" ></kendo-grid-column>
      <kendo-grid-column :field="'Status'"
                         :title="'STATUS'" ></kendo-grid-column>
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
        },
    data: function() {
          return { 
            pageable: {
                        pageSizes: [5, 10, 20],
                        buttonCount: 5,
                        pageSize: 5
                      },

            receive:[{ EmployeeId:'',
                       Gender: '',
                       LeaveType:'',
                       Leavebalance: '',
                       FromDate: '',
                       ToDate: '', 
                       NoofDays: '',
                       Status: ''
                    }],
                }
            },
            created() {
                this.get()
            },

            methods: {
                get(){
                    axios.get('http://localhost:8001/request', {})
                    .then((response) => {
                    console.log(response.data)
                    this.receive = response.data
                    this.get
                })
                    .catch((error) => {
                    console.log(error)
                })
                console.log("uuupppdddaaattteee", this.receive)
                },
                
            }
        };

</script>