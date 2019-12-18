$(document).ready(function(){
    // Schedule at datetime picker
    var $statusInput = $('#inputStatus');
    var $scheduledAtInput = $('#inputScheduledAt');
    var $scheduledAtContainer = $('#input-scheduled-at-container');
    if($scheduledAtInput.length > 0 && $statusInput.length > 0) {
        // Check for the value, if it's scheduled at, then show the scheduled at datetimepicker
        if($statusInput.val() == 1) {
           $scheduledAtContainer.show(); 
        } else {
            $scheduledAtContainer.hide();
        }
        
        // Hide or display the datetime picker
        $statusInput.change(function(){
            if($(this).val() == 1) {
                $scheduledAtContainer.show(); 
             } else {
                 $scheduledAtContainer.hide();
             }
        });

        $scheduledAtInput.datetimepicker({
            format: 'DD MMM YYYY h:mm A',
            stepping: 5,
            timeZone: "Asia/Kuala_Lumpur"
        });
    }

    let csrfToken = document.getElementsByName("gorilla.csrf.Token")[0].value;
    // Initialize WYSIWYG editor
    tinymce.init({
        selector: 'textarea#inputContent',
        height: 500,
        plugins: "image table",
        images_upload_handler: function (blobInfo, success, failure) {
            var formData = new FormData();
            formData.append('image', blobInfo.blob());
            // To upload the file to the system first
            axios({
                method: "post",
                url: "/dashboard/post/upload",
                data: formData,
                headers: {"X-CSRF-Token": csrfToken, "Content-Type": "multipart/form-data"},
            })
            .then(function (response) {
                // handle success
                data = response["data"];
                if(data.success){
                    success(data.data);
                } else {
                    failure(data.message);
                }                        
            })
            .catch(function (error) {
                // handle error
                failure(error);
            });
        },
        file_picker_callback: function(callback, value, meta) {
            // Provide file and text for the link dialog
            if (meta.filetype == 'file') {
              callback('mypage.html', {text: 'My text'});
            }
        
            // Provide image and alt text for the image dialog
            if (meta.filetype == 'image') {
              callback('myimage.jpg', {alt: 'My alt text'});
            }
        
            // Provide alternative source and posted for the media dialog
            if (meta.filetype == 'media') {
              callback('movie.mp4', {source2: 'alt.ogg', poster: 'image.jpg'});
            }
          }
    });

    // Validate the create/edit post form
    $.validator.addMethod("laterThanNow",function(value, element){
        var now = new moment().add(15, "minutes");
        var timeSelected = moment($('#inputScheduledAt').val());
        return this.optional(element) || timeSelected.isAfter(now);
    }, "The schedule datetime must be at least 15 minutes later.");

    $("#create-edit-post-form").validate({
        rules: {
            title: {
                required: true
            },
            content: {
                required: true
            },
            scheduled_at: {
                date: true,
                required: function(element){
                    return $("#inputStatus").val() == 1;
                },
                laterThanNow: true
            }
        },
        messages: {
            title: {
                required: "Title is a mandatory field."
            },
            content: {
                required: "Content is a mandatory field."
            },
            scheduled_at: {
                date: "Invalid datetime format.",
                required: "Schedule datetime must be set.",
                laterthanNow: "The schedule datetime must be at least 15 minutes later."
            }
        },
        submitHandler: function(form) {
            form.submit();            
            toggleLoading();
        }
    });

    var $btnDeletePost = $('#btn-delete-post');
    var $deletePostForm = $('#delete-post-form');
    if($btnDeletePost.length > 0 && $deletePostForm) {
        $btnDeletePost.click(function(e){
            e.preventDefault();
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
                    $deletePostForm.submit();
                    toggleLoading();
                }
            })
        });
    }
});