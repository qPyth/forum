$(document).ready(function(){
    $("#signin-button").click(function(){
        $("#loginModal").css("display", "block");
    });
    $("#registerBtn").click(function(){
        $("#loginModal").css("display", "none");
        $("#registerModal").css("display", "block");
    });
    $("#closeLoginModal").click(function(){
        $("#loginModal").hide();
    });
    $("#closeRegisterModal").click(function(){
        $("#registerModal").hide();
    });

    $("#login-submit").click(function(){
        event.preventDefault();
        let formData = $("#signin-form").serialize();
        $.ajax({
            url: "/sign-in",
            method: "post",
            dataType: "json",
            data: formData,
            success: function(data, textStatus, jqXHR){
                console.log("succ")
                console.log(jqXHR.responseText)
                let response = JSON.parse(jqXHR.responseText);
                console.log(response);
                window.location.href = response.referer;
            },
            error: function(xhr, textStatus, errorThrown) {
                let errorResponse = JSON.parse(xhr.responseText);
                $("#error-signin").text(errorResponse.errorText);
            }
        });
    });

    $("#signup-submit").click(function(){
        event.preventDefault();
        let formData = $("#signup-form").serialize();
        $.ajax({
            url: "/sign-up",
            method: "post",
            dataType: "json",
            data: formData,
            success: function(){
                console.log("succ")
                window.location.href = "http://localhost:8080/";
            },
            error: function(xhr, textStatus, errorThrown) {
                let errorResponse = JSON.parse(xhr.responseText);
                $("#registerModal").css({"height": "450px"});
                $("#error-signup").text(errorResponse.errorText);
            }
        });
    });

    /*-------------------USERPROFILE-----------------------*/

    $("#user-profile-button").click(function () {
        event.preventDefault()
        $("#profile-modal").css("display", "block");
    })

    $('body').click(function (event)
    {
        if(!$(event.target).closest('#user-profile-button').length && !$(event.target).is('#user-profile-button')) {
            $(".profile-modal").hide();
        }
    });

    /*LIKES AND DISLIKES*/


    var likeCount = 0;
    var dislikeCount = 0;

    $('.like, .dislike').click(function() {
        var isPost = $(this).closest('.post-section').length > 0;
        var isComment = $(this).closest('.comment').length > 0;
        var id = $(this).closest('[data-id]').data('id');
        var action = $(this).hasClass('like') ? 'like' : 'dislike';
        var data = {
            'id': id,
            'action': action,
            'is_post': isPost,
            'is_comment': isComment
        }
        console.log(data)
        $.ajax({
            url: "/vote",
            method: "post",
            data: JSON.stringify(data),
            contentType: "application/json",
            success: function(data, textStatus, jqXHR){
                console.log(data)
                var likeCountElement = $(this).closest('.like-outer').find('.likeCount');
                var dislikeCountElement = $(this).closest('.like-outer').find('.dislikeCount');
                likeCountElement.text(data.like_count);
                dislikeCountElement.text(data.dislike_count);

                $(this).closest('.like-outer').find('i').removeClass('liked disliked');


                if (data.action === 'like') {
                    $(this).addClass('liked');
                } else if (data.action === 'dislike') {
                    $(this).addClass('disliked');
                }
            }.bind(this),  // Привязка контекста this к функции success
            error: function(jqXHR, textStatus, errorThrown) {
                if (jqXHR.status===403)  {
                    $("#loginModal").css("display", "block");
                }
            }
        });
    });
});





