$(document).ready(function(){    
    var $form = document.getElementById("create-edit-company-form");
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
            console.log(response);
            $form.reset();

            // Set the value in the dialog
            $form.action = editURL;
            data = response["data"];
            
            // Unhide the loading
            toggleLoading();   
            console.log(data);
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
            console.log(error);
            toggleLoading();
            Swal.fire({
                type: 'error',
                title: 'Oops...',
                text: 'Something went wrong!',
            })
        });
    })
});
