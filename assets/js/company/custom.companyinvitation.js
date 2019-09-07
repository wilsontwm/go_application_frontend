$(document).ready(function(){    
    $invitationForm = document.getElementById("respond-company-invitation-form");
    $(".invitation-link").click(function(e) {
        e.preventDefault();
        var $row = $(this).closest(".invitation-row");
        if($row.length > 0) {
            var $btn = $row.find(".btn-view-invitation");
            if($btn.length > 0) {
                $btn.click();
            }
        }
    });

    $(".btn-view-invitation").click(function(){
        var id = $(this).data('id');
        var companyName = $(this).data('company');
        var status = $(this).data('status');
        var message = $(this).data('message');
        var senderName = $(this).data('sender');
        var senderEmail = $(this).data('senderemail');

        $invitationForm.action = "/dashboard/invite/incoming/" + id + "/respond";
        $("#company-name-container").html(companyName);
        if(message.length > 0) {
            var messageString = '<figure class="quote"><div class="quote-body">' + message + '</div></figure>';
            $("#message-container").html(messageString);
        } else {
            $("#message-container").html('');
        }
        $("#sender-name-container").html(senderName);
        var senderEmailString = '<a href="mailto:' + senderEmail + '">' + senderEmail + '</a>';
        $("#sender-email-container").html(senderEmailString);

        if(parseInt(status) == 0) {
            $("#btn-accept-invitation").attr("disabled", false);
            $("#btn-decline-invitation").attr("disabled", false);
        } else {
            $("#btn-accept-invitation").attr("disabled", true);
            $("#btn-decline-invitation").attr("disabled", true);
        }
        $("#modal-view-invitation").modal("show");
    });
    
    $("#btn-accept-invitation").click(function(e) {
        e.preventDefault();
        $("#inputResponse").val("accept");
        $invitationForm.submit();
        $("#modal-view-invitation").modal("hide");
        toggleLoading();
    });

    $("#btn-decline-invitation").click(function(e) {
        e.preventDefault();
        $("#inputResponse").val("decline");
        $invitationForm.submit();
        $("#modal-view-invitation").modal("hide");
        toggleLoading();
    });

});
