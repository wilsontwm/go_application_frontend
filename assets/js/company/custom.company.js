$(document).ready(function(){    
    var $form = document.getElementById("create-edit-company-form");
    var $deleteform = document.getElementById("delete-company-form");

    // Create / edit the company
    $("#btn-create-company").click(function(){
        $("#modal-create-edit-company-title").html("Create new company");
        $form.action = $(this).data('url');
        $form.reset();
        

        $("#modal-create-edit-company").modal("show");
    });

    $("#btn-edit-company").click(function(){
        $("#modal-create-edit-company-title").html("Edit company");
        $form.action = $(this).data('url');

        $("#modal-create-edit-company").modal("show");
    });

    // Get the company detail when click on edit button
    $(".btn-edit-company").click(function(){
        // Show loading first
        toggleLoading();
        var id = $(this).data("id");
        var url = "/dashboard/company/" + id + "/show/json";
        var editURL = "/dashboard/company/" + id + "/update";
        var data;
        axios.get(url)
        .then(function (response) {
            // handle success
            $form.reset();

            // Set the value in the dialog
            $form.action = editURL;
            data = response["data"];
            
            // Unhide the loading
            toggleLoading();   
            
            if(data["data"] != null) {
                $("#modal-create-edit-company-title").html(data["data"]["Name"]);
                $form.elements.namedItem("name").value = data["data"]["Name"];
                $form.elements.namedItem("slug").value = data["data"]["Slug"];
                $form.elements.namedItem("description").value = data["data"]["Description"];
                $form.elements.namedItem("email").value = data["data"]["Email"];
                $form.elements.namedItem("phone").value = data["data"]["Phone"];
                $form.elements.namedItem("fax").value = data["data"]["Fax"];
                $form.elements.namedItem("address").value = data["data"]["Address"];
                $("#modal-create-edit-company").modal("show");
            } else {
                Swal.fire({
                    type: 'error',
                    title: 'Oops...',
                    text: 'Something went wrong! Please refresh the page.',
                })
            }
                     
        })
        .catch(function (error) {
            // handle error
            toggleLoading();
            Swal.fire({
                type: 'error',
                title: 'Oops...',
                text: 'Something went wrong!',
            })
        });
    });    

    $("#create-edit-company-form").validate({
        rules: {
            name: {
                required: true
            },
            slug: {
                required: true
            }
        },
        messages: {
            name: {
                required: "Name is a mandatory field."
            },
            slug: {
                required: "Slug is a mandatory field."
            }
        },
        submitHandler: function(form) {
            form.submit();    
            $("#modal-create-edit-company").modal("hide");        
            toggleLoading();
        }
    });

    // Set the URL for delete when click on delete button
    $(".btn-delete-company").click(function(){        
        var id = $(this).data("id");
        var deleteURL = "/dashboard/company/" + id + "/delete";
        deleteCompany(deleteURL);
    });

    $("#btn-delete-company").click(function(){
        var deleteURL  = $(this).data('url');
        deleteCompany(deleteURL);
    });

    function deleteCompany(url) {
        $deleteform.action = url;
        Swal.fire({
            title: 'Are you sure?',
            text: "All the relevant data will be deleted!",
            type: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, delete it!'
          }).then((result) => {
            if (result.value) {   
                $deleteform.submit();
                toggleLoading();
            }
          })
    }
});


