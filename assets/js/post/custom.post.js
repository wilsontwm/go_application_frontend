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

    // Initialize WYSIWYG editor
    tinymce.init({
        selector: 'textarea#inputContent',
        height: 500,
        plugins: "image table",
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
});