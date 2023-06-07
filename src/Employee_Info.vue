<template>
<br>
<div class="head">
  <h4> EMPLOYEES LEAVE BALANCE DETAILS </h4>
</div>
<hr>
<div class="leavebtn">
    <button id="btn" @click="click()">APPLY FOR LEAVE <i class="bi bi-pencil-square"></i> </button>
</div>
<div v-bind:class="{active:isactive}" class="modal">
<div class="modal-content">
<div class="button">  
    <h1 class="modal-title fs-5" > APPLY LEAVE </h1>
</div>
<hr>
<form>
  <table class="tablecontent" >
    <tr>
        <td><label> EMPLOYEE ID </label></td>
        <td class="idata"> <kendo-dropdownlist class="inputf"
                            v-model="send.EmployeeId"
                            :data-source="getempid"
                            :option-label="'Select Empid'" /></td>
    </tr>
    <tr>
      <td><label> LEAVE TYPE </label></td>
      <td class="idata"> <kendo-dropdownlist class="inputf"
                          @close="getinfo()" 
                          v-model="send.LeaveType"
                          :data-source="leavetype"
                          :option-label="'Select Leavetype'"/></td>
    </tr>
    <tr>
        <td><label> GENDER </label></td>
        <td class="idata"> <input class="inputf"
                            v-model="GenderData" type="text" placeholder="Gender" disabled/></td>
        <!-- <td class="idata"> <p> {{ GenderData }} </p> </td> -->
    </tr>
    <tr>
        <td><label> LEAVE BALANCE </label></td>
        <td class="idata"> <input class="inputf"
                            v-model="LeavebalanceData" type="text" placeholder="Balance Leave" disabled/></td>
        <!-- <td class="idata"><p> {{ LeavebalanceData }} </p></td> -->
    </tr>
    <tr>
        <td><label> FROM DATE </label></td>
        <td class="idata"> <datepicker class="inputf" 
                            :format="'dd/MM/yyyy'" 
                            v-model="send.FromDate"/> </td>
    </tr>
    <tr>
        <td><label> TO DATE </label></td>
        <td class="idata"> <datepicker class="inputf" 
                            :format="'dd/MM/yyyy'" 
                            @change="daysDifference" 
                            v-model="send.ToDate"/> </td>
        <!-- <td> <DateRangePicker :default-value="defaultValue" @change="daterange" /> </td> -->
    </tr>
    <tr>
        <td> <label> NO OF DAYS </label></td>
        <td class="idata"> <input class="inputf" 
                            v-model="send.NoofDays" type="number" placeholder="No of Days" disabled/></td>
        <!-- <td class="idata"><p> {{ this.send.NoofDays }} </p></td> -->
      </tr>
    <tr>
        <td class="hide"><KInput :style="{ color: textColor }" v-model="send.Status" type="text" ></KInput></td>
    </tr>
  </table>
      <hr>
    <div class="button" >
        <button class="btn btn-success" @click="submitform" > SUBMIT <i class="bi bi-download"></i> </button>
        <button class="btn btn-danger" id="cancel" @click="close" > CANCEL <i class="bi bi-x-lg"></i> </button>
    </div>
</form>
</div>
</div>
<br> 
<div id="kgridheader">
      <kendo-grid :data-source="receive"
                  :pageable="pageable"
                  :filterable="true"
                  :noRecords="true"
                  :height="266">
                    <kendo-grid-column :field="'Empid'" :title="'EMPLOYEE ID'" ></kendo-grid-column>
                    <kendo-grid-column :field="'Name'" :title="'NAME'"></kendo-grid-column>
                    <kendo-grid-column :field="'Gender'" :title="'GENDER'"></kendo-grid-column>
                    <kendo-grid-column :field="'CL'" :title="'CASUAL LEAVE BALANCE '"></kendo-grid-column>
                    <kendo-grid-column :field="'SL'" :title="'SICK LEAVE BALANCE'"></kendo-grid-column>
                    <kendo-grid-column :field="'PL'" :title="'PATERNITY LEAVE BALANCE'"></kendo-grid-column>
      </kendo-grid>
</div>

<!-- Vue-Loading overlay animation -->
<div>  
      <loading 
        :active='isLoading' 
        :is-full-page="fullPage" 
        :loader='loader'
        :color='color'
        :width='width' />
</div>

<div>
    <!-- <DateRangePicker :default-value="defaultValue" @change="daterange" /> -->
    <!-- <p> Start date : {{startdate}} </p>
    <p> End date : {{enddate}} </p> -->
</div>

</template>

<script>
        import { Input } from '@progress/kendo-vue-inputs';
        import { DatePicker } from '@progress/kendo-vue-dateinputs';
        import { DropDownList } from '@progress/kendo-dropdowns-vue-wrapper';
        import { Grid, GridColumn } from '@progress/kendo-grid-vue-wrapper';
        import 'vue-loading-overlay/dist/css/index.css';
        import Loading from 'vue-loading-overlay';
        import axios from 'axios';
        import Swal from 'sweetalert2';
        // import { DateRangePicker } from "@progress/kendo-vue-dateinputs";
        import 'bootstrap/dist/css/bootstrap.css'
        import 'bootstrap/dist/js/bootstrap.js'


export default {
    components: {  
        'kendo-grid'          : Grid,
        'KInput'              : Input,
        'kendo-grid-column'   : GridColumn,
        'kendo-dropdownlist'  : DropDownList,
        'datepicker'          : DatePicker,
        Loading,
        // DateRangePicker,
        },

    data: function(){
          return { 
            ismodal: false,
            leavetype: ['Casual Leave','Sick Leave','Paternity Leave'],
            isactive: false,
            getempid: [],
            GenderData: "",
            LeavebalanceData: "",
            isLoading: false,
            fullPage: true,
            loader: 'dots',
            color: 'rgb(17, 107, 139)',
            width: '60px',
            textColor: 'red',

            // startdate: "",
            // enddate: "",

    pageable: { pageSizes: [1, 5, 10, 20],
                buttonCount: 5,
                pageSize: 5,
              },

    send:{ EmployeeId: '',
            LeaveType: '',
            FromDate: '',
            ToDate: '',
            NoofDays: '',
            Status: 'Pending!',
          },

    receive:{ Empid: '',
              Name: '',
              Gender: '',
              CL: '',
              SL: '',
              PL: '',
            },

    receivedata:[{ Empid: '',
                   Name: '',
                   Gender: '',
                   CL: '',
                   SL: '',
                   PL: '',
                }],
              }
            },

    created() {
          this.get()
        },

    methods: {

        // daterange(e){

        //     console.log("Event :",e)

        //     console.log("Start Date :",e.value.start.getDate()+"/"+ (e.value.start.getMonth()+1) +"/"+ e.value.start.getFullYear())
        //     console.log("End Date :",e.value.end.getDate()+"/"+ (e.value.end.getMonth()+1) +"/"+ e.value.end.getFullYear())

        //     this.startdate = e.value.start.getDate()+"/"+ (e.value.start.getMonth()+1) +"/"+ e.value.start.getFullYear()
        //     this.enddate = e.value.end.getDate()+"/"+ (e.value.end.getMonth()+1) +"/"+ e.value.end.getFullYear()

        //   },
        headerCellProps: function() {
      // Set the background color of the column header to red
      return { style: { backgroundColor: "red" } };
    },

        daysDifference() {  
            var dateI1 = this.send.FromDate;  
            var dateI2 = this.send.ToDate;  
            var date1 = new Date(dateI1);  
            var date2 = new Date(dateI2);  
            var time_difference = date2.getTime() - date1.getTime();  
            var result = time_difference / (1000 * 60 * 60 * 24);  
            this.send.NoofDays = result+1
          },

        click(){
                this.isLoading = true;
                setTimeout(() => {
                this.isLoading = false}, 400)
                this.isactive=true
                this.getid()
              },

        close(){
                this.isactive=false
                this.clear()
              },

        clear(){
                this.send.EmployeeId = [] ,
                this.send.LeaveType = [] ,
                this.GenderData = "" ,
                this.LeavebalanceData = "" ,
                this.send.FromDate = "" ,
                this.send.ToDate = "" ,
                this.send.NoofDays = ""
              },
              
        get(){
              axios.get('http://localhost:8007/request', {})
              .then((response) => {
              console.log(response.data)
              this.receive = response.data
              this.get
              })
              .catch((error) => {
              console.log(error)
              })
            },

        submitform(){
            
          if ( this.send.EmployeeId == "" ||
                this.send.LeaveType == "" ||
                this.send.FromDate == "" ||
                this.send.ToDate == "" ||
                this.send.NoofDays == "" )
                {Swal.fire({
                    icon: 'error',
                    title: "Can't Submit",
                    text: 'Fill all the fields!',
                  })}
                // {alert("Fill all the fields")}
          else if (this.GenderData == "Male" && this.send.LeaveType == "Paternity Leave")
                {Swal.fire({
                    icon: 'warning',
                    title: "Can't Submit",
                    text: 'You are not eligible to apply for paternity leave!',
                  })}
          else if (this.LeavebalanceData < this.send.NoofDays)
                {Swal.fire({
                    icon: 'warning',
                    title: "Can't Submit",
                    text: 'Your no of days leave taken is more than leave balance so you are not able to apply for leave!',
                  })}
          else if (this.send.NoofDays <= 0)
                {Swal.fire({
                    icon: 'warning',
                    title: "Can't Submit",
                    text: 'Select a valid date!',
                  })}
          else  {axios.post("http://localhost:8001/insert", {
                        'EmployeeId'    : this.send.EmployeeId,
                        'Gender'        : this.GenderData,
                        'LeaveType'     : this.send.LeaveType,
                        'Leavebalance'  : this.LeavebalanceData,
                        'FromDate'      : this.send.FromDate,
                        'ToDate'        : this.send.ToDate,
                        'NoofDays'      : this.send.NoofDays,
                        'Status'        : this.send.Status,
                      })
                        .then(response => {
                        console.log(response);
                        this.close()
                        this.clear()
                      })
                        .catch(e => {
                        this.errors.push(e)
                      });
                      Swal.fire(
                            'Success!',
                            'Leave applied successfully...!',
                            'success'
                       )
                       console.log("ffffrrrooommm", this.send.FromDate)
                      console.log("tttttooooo", this.send.ToDate )
                      // const Toast = Swal.mixin({
                      //     toast: true,
                      //     position: 'top',
                      //     showConfirmButton: false,
                      //     timer: 3000,
                      //     timerProgressBar: true,
                      //     didOpen: (toast) => {
                      //       toast.addEventListener('mouseenter', Swal.stopTimer)
                      //       toast.addEventListener('mouseleave', Swal.resumeTimer)
                      //     }
                      //   })
                      //   Toast.fire({
                      //     icon: 'success',
                      //     title: 'Leave Applied Successfully!'
                      //   })
                    }
                  },

        getinfo(){
            axios.post("http://localhost:8007/getempid", {
                'Empid': this.send.EmployeeId })
                .then((response) => {
                this.receivedata = response.data
                console.log("dddddd", this.receivedata)
                this.GenderData = this.receivedata[0][2].Value
                  if( this.send.LeaveType == 'Casual Leave'){
                    this.LeavebalanceData = this.receivedata[0][3].Value
                  }else if(this.send.LeaveType == 'Sick Leave'){
                    this.LeavebalanceData = this.receivedata[0][4].Value
                  }else{
                    this.LeavebalanceData = this.receivedata[0][5].Value
                  }
                })
                  .catch(e => {
                  this.errors.push(e)
                });
              },

        getid(){
            var id = this.receive
            var Id=[]
            for (var i=0; i<id.length; i++ ){
                Id.push(this.receive[i].Empid)
            }
            this.getempid = Id
            console.log('###',this.getempid )
            },    
          }
        };
</script>

<style>

  #kgridheader{
    background:  #16aecd!important;
    color: white!important;
    font-weight: bold !important;
  }
  .idata{
    padding-left: 40px;
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
  .modal-content {
    background-color: #fefefe;
    margin: auto;
    padding: 20px;
    border: 1px solid rgb(128, 127, 127);
    width:520px;
    height: 600px;
    border-radius: 15px;
  }
  #btn{
    position: relative;  left:10px; 
    color: #ffffff;
    background-color: #11abca;
    border: 2px solid rgb(17, 107, 139);
    border-radius: 8px;
    padding: 5px 18px;
    cursor: pointer;
    transition: all 0.4s ease-in-out;
  }
  #btn:hover{
    transform: scale(1.1);
    color: #ffffff;
    border-color: #000;
    background-color: #11abca;
    border: 2px solid rgb(14, 88, 107);
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
  .inputf{
    width: 200px;
    height: 38px;
  }
  .close:hover,.close:focus {
    color: #000;
    text-decoration: none;
    cursor: pointer;
  }
  .hide{
    display: none;
  }
  .active{
    display:block;
  }
  .tablecontent{
    padding: 25px;
  }
  .button{
    text-align: center;
    margin-top: 15px;
  }
  td,tr{
    padding: 10px;
  }
  #cancel{
    margin-left: 20px;
  }
  input{
    border: none;
    border-color: transparent;
    background-color: transparent;
  }

</style>
