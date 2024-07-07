$(document).ready(function () {

    $('#files').on("change", function() {
        let val = $(this).val(), btn = $('#submit');
        val ? btn.removeAttr("disabled") : btn.attr("disabled");
    });

    $('#uploadForm').on('submit', function (event) {
       event.preventDefault();

        var formData = new FormData(this);
        $.ajax({
            xhr: function () {
                var xhr = new window.XMLHttpRequest();
                xhr.upload.addEventListener('progress', function (event) {
                    if (event.lengthComputable) {
                        var percentComplete = Math.round((event.loaded / event.total) * 100);
                        $('#progressBar').width(percentComplete + '%');
                        $('#progressBar').attr('aria-valuenow', percentComplete);
                        $('#progressBar').text(percentComplete + '%');
                    }
                }, false);
                return xhr;
            },
            url: '/upload',
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            beforeSend: function () {
                $('.progress').show();
                $('#progressBar').width('0%');
                $('#progressBar').attr('aria-valuenow', '0');
                $('#progressBar').text('0%');
            },
            success: function (response) {
                // Redirect to the desired URL after successful upload
                window.location.href = '/download';

                $('#uploadForm').trigger("reset");
                $('.progress').hide();
                $('#submit').attr('disabled', 'disabled');
            },
            error: function () {
                alert('An error occurred while uploading the files.');
                $('.progress').hide();
            }
        });
    });
});