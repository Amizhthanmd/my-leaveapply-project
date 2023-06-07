<template>
<body class="backgr">
<br>
<br>
<div class="formhead">
    <h4> ADD NEW ANNOUNCEMENT </h4>
</div>
<div class="container">
    <form  enctype="multipart/form-data" action="http://localhost:8007/getemail" method="post" validation-schema="schema" >
        <table>
            <tr>
                <td><label> Audience </label></td>
                 <td class="td"><kendo-dropdownlist style="width:25%;" class="inputw" 
                    v-model="audience"
                    name="audience"
                    :data-items="dataCategories"
                    @change="categoryChange"/></td>
            </tr>
            <tr>
                <td><label> Select type </label></td>
                <td class="td"><kendo-dropdownlist style="width:25%" class="inputw"
                    v-model="selecttypes"
                    name="selecttype" 
                    :data-items="types"/></td>
            </tr>
            <tr>
                <td><label> Title </label></td>
                <td class="td" ><KInput required style="width:25%" class="inputw" v-model="title"  name="title" placeholder="Enter Title" type="text" ></KInput></td>
            </tr>
            <tr>
                <td><label> CC </label></td>
                <td class="td"><KInput style="width:25%" class="inputw" v-model="cc" name="cc" placeholder="CC" type="text" ></KInput></td>
            </tr>
            <tr>
                <td><label> Body </label></td>
                    <td style="padding-right: 20%;"> <froala  v-model="body" id="edit" :tag="'textarea'" :config="config" name="body"></froala></td>
            </tr>
            <tr>
                <td><label> Attachment </label></td>
                <td class="td"><KInput style="width:25%" class="inputw" v-model="attachments" name="myFile" type="file" ></KInput></td>
            </tr>
            <tr>
                <td></td>
                <td class="td"> <button @click.prevent="clear" class="btn btn-dark" > <i class="bi bi-arrow-clockwise"></i> Reset </button> <button id="postbtn"  class="btn btn-success" > Post <i class="bi bi-send"></i> </button></td>
            </tr>
        </table>
    </form>
</div>
</body>
</template>

<script>

        import { Input } from '@progress/kendo-vue-inputs';
        import { DropDownList } from '@progress/kendo-vue-dropdowns';

    export default {
        components: {  
            'kendo-dropdownlist'  : DropDownList,
            'KInput'              : Input,
        },
        data: function(){
            return{
                config: {
                events: {
                    initialized: function () {
                    console.log('initialized')
                }
            }
        },
                dataCategories : ["Gender","Skills"],
                category: null,
                product: null,
                audience : [],
                selecttypes : [],
                types: [],
                title : "",
                cc : "",
                body: "",
                file: "",
                attachments: "",
                }
            },

        methods: {

            categoryChange(e) {

                if(e.value == "Gender"){
                    this.types = ["Male","Female"]
                }else if(e.value == "Skills"){
                    this.types = ["Golang","Python","Java","Angular","Reactjs","Vuejs","MongoDB"]
                }
            },
            
        clear(){
                this.audience = [],
                this.selecttypes = [],
                this.title = '',
                this.cc = '',
                this.body = '',
                this.attachments = ''
            },
        }
    }
            // getemailid(){

            //     let formData = new FormData();
            //     console.log("ffffffffffffffff :",formData);
            //     formData.append('file', this.file);

            //     const current = new Date();
            //     const date = `${current.getDate()}/${current.getMonth()+1}/${current.getFullYear()}`;
            //     console.log(this.audience.categoryName)
            //         axios.post("http://localhost:8007/getemail", {
            //             'Audience' : this.audience.categoryName ,
            //             'Selecttype' : this.selecttype.productName,
            //             'Title' : this.title,
            //             'Postedon' : date,
            //             formData,
            //             headers: {
            //             'Content-Type': 'multipart/form-data'
            //             }
            //         })
            //         .then((response) => {
            //             this.receive = response.data
            //            console.log("mailssss",this.receive)
            //         })
            //         .catch(e => {
            //             this.errors.push(e)
            //         });
            //     },

        

</script>

<style>
.container{
   padding-left: 250px;
}

.inputw{
    height: 36px;
}
.backgr{
    background-color: rgb(137, 179, 194);
}
.td{
    padding: 18px;
}
label{
    font-size: larger;
}
#postbtn{
    margin-left: 15px;
}
.formhead{
    margin-left: 40%;
    color: rgb(37, 37, 37);
}
</style>


