<template>
<br>
<div class="head">
    <h4> PENDING FOR APPROVAL </h4>
</div>
<hr>
<div v-bind:class="{active:isactive}" class="modal">
<div class="modal-contents">
<div class="apphead">
    <h1 class="modal-title fs-5" > APPROVAL FORM </h1>
</div>
<hr>
<form>
  <table class="tablecontent" >
    <tr>
        <td><label> EMPLOYEE ID </label></td>
        <td class="data"> {{EmployeeId}} </td>
    </tr>
    <tr>
        <td><label> GENDER </label></td>
        <td class="data"> {{Gender}} </td>
    </tr>
    <tr>
        <td><label> LEAVE TYPE </label></td>
        <td class="data"> {{LeaveType}} </td>
    </tr>
    <tr>
        <td><label> LEAVE BALANCE </label></td>
        <td class="data"> <input v-model="Leavebalance" disabled/>  </td>
    </tr>
    <tr>
        <td><label> FROM DATE </label></td>
        <td class="data"> {{FromDate}} </td>
    </tr>
    <tr>
        <td><label> TO DATE </label></td>
        <td class="data" > {{ToDate}} </td>
    </tr>
    <tr>
        <td><label> NO OF DAYS </label></td>
        <td class="data"> <input v-model="NoofDays" disabled/>  </td>
    </tr>
    <tr>
        <td><label> STATUS </label></td>
        <td class="data"><kendo-dropdownlist class="input"
                          :data-source="status" 
                          v-model="Status"/></td>
    </tr>
  </table>
<hr>
<div class="btnbtn">
      <button class="btn btn-success"  v-on:click="approveleave(); Delete()"> SUBMIT <i class="bi bi-download"></i> </button>
      <button class="btn btn-danger" id="subbtn"  @click="close" > CANCEL <i class="bi bi-x-lg"></i> </button>
</div>
</form>
</div>
</div>
  
<div id="kgridheader">
      <kendo-grid :data-source="receive" ref="maingrid"
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
      <kendo-grid-column :command="[{name:'View',click: view}]"
                         :title="'APPROVE'" ></kendo-grid-column>
      </kendo-grid>
</div>
<div class="note">
    <span> Note: This grid shows only Pending! details </span>
</div>
</template>

<script>

      import { DropDownList } from '@progress/kendo-dropdowns-vue-wrapper';
      import { Grid, GridColumn } from '@progress/kendo-grid-vue-wrapper';
      import axios from 'axios';
      import $ from 'jquery';
      import Swal from 'sweetalert2'

export default {
    components: {  
        'kendo-grid'          : Grid,
        'kendo-dropdownlist'  : DropDownList,
        'kendo-grid-column'   : GridColumn,
        },
    data: function() {
          return { 
            pageable: {
                        pageSizes: [5, 10, 20],
                        buttonCount: 5,
                        pageSize: 5
                      },
            isactive:false,
            status: ['Pending!','Approved!','Rejected!'],

                      EmployeeId:'',
                      Gender: '',
                      LeaveType:'',
                      Leavebalance: '',
                      FromDate: '',
                      ToDate: '', 
                      NoofDays: '',
                      Status: '',

            receive:{ EmployeeId:'',
                      Gender: '',
                      LeaveType:'',
                      Leavebalance: '',
                      FromDate: '',
                      ToDate: '', 
                      NoofDays: '',
                      Status: ''
                   },
                }
            },
            created() {
                this.get()
                },
            methods: {
                close(){
                this.isactive=false
            },

            approveleave(){
                axios.post("http://localhost:8001/update", {
                        EmployeeId:   this.EmployeeId,
                        Gender:       this.Gender,
                        LeaveType:    this.LeaveType,
                        Leavebalance: this.Leavebalance,
                        FromDate:     this.FromDate,
                        ToDate:       this.ToDate,
                        NoofDays:     this.NoofDays,
                        Status:       this.Status,
                        })
                        .then(response => {
                        console.log(response)
                        this.get()
                        this.close()
                        })
                        .catch((error) => {
                        console.log('##',error)
                        });
                        Swal.fire(
                            'Success!',
                            'Leave status updated...!',
                            'success'
                          )
                      },
            get(){
                  axios.get('http://localhost:8001/sendtoapproval', {})
                      .then((response) => {
                      console.log('@##@',response.data)
                      this.receive = response.data
                      })
                      .catch((error) => {
                      console.log(error)
                      })
                 },
                view(ev){
                    var gridWidget = this.$refs.maingrid.kendoWidget(); 
                    var tr = $(ev.target).closest("tr");
                    gridWidget.dataItem(tr);
                    this.EmployeeId = gridWidget.dataItem(tr).EmployeeId,
                    this.Gender = gridWidget.dataItem(tr).Gender,
                    this.LeaveType = gridWidget.dataItem(tr).LeaveType,
                    this.Leavebalance = gridWidget.dataItem(tr).Leavebalance,
                    this.FromDate = gridWidget.dataItem(tr).FromDate,
                    this.ToDate = gridWidget.dataItem(tr).ToDate,
                    this.NoofDays = gridWidget.dataItem(tr).NoofDays,
                    this.Status = gridWidget.dataItem(tr).Status,
                    console.log('aaa',this.Leavebalance)
                    this.isactive=true
                },
                Delete(){
                  axios.post("http://localhost:8001/delete",{ 
                        EmployeeId:  this.EmployeeId,
                        Status:      this.Status,})
                        .then(response => {
                        console.log(response)
                        })
                        .catch(e => {
                        this.errors.push(e)
                        })
                        this.get()
                        }
                      }
                    };
</script>

<style>

  .data{
    padding-left: 80px;
  }
  .modal {
    display: none; 
    position: fixed;
    z-index: 9; 
    padding-top: 100px; 
    left: 0;
    top: 0;
    width: 100%; 
    height: 100%; 
    overflow: auto; 
    background-color: rgb(0,0,0);
    background-color: rgba(0,0,0,0.4); 
  }
  .modal-contents {
    background-color: #fefefe;
    margin: auto;
    padding: 20px;
    border: 1px solid rgb(128, 127, 127);
    width:500px;
    height: 615px;
    border-radius: 15px;
  }
  .btnbtn{
    text-align: center;
  }
  #subbtn{
    margin-left: 20px;
  }
  .leavebtn{
    text-align: center;
  }
  .head{
    text-align: center;
  }
  .close {
    color: #aaaaaa;
    float: right;
    font-size: 35px;
    font-weight: bold;
    text-align: right;
  }
  .input{
    width: 300px;
  }
  .close:hover,.close:focus {
    color: #000;
    text-decoration: none;
    cursor: pointer;
  }
  .active{
    display:block;
  }
  .apphead{
    text-align: center;
  }
  label {
    font-weight: 550;
  }
  input{
    border: none;
    border-color: transparent;
    background-color: transparent;
  }
  .note{
    position: relative;  left:10px; 
    font-size: smaller;
  }

</style>
    



