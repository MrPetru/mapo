

//$(document).ready(function() { formNewUser(); });

function formNewUser() {
    form = $("<form>", {
                action: "/admin/user",
                method: "POST",
                id: "newuserform"
    }
    );

    var formTitle = $("<h1/>");
    formTitle.text("Create New User");
    form.append(formTitle);

    form.append(text("Login:", "login"));
    form.append(text("Password:", "password"));
    form.append(text("Displa Name:", "name"));
    form.append(textarea("Description:", "description"));
    form.append(text("Email:", "email"));

    var errorContainer = $("<div/>", {name: "error:"+name});
    errorContainer.addClass("error");
    form.append(errorContainer);

    form.append(submit());

    form.addClass("formcontainer");
    form.appendTo("body");
    };

function submit() {
    var d = $("<div/>");
    d.addClass("inputBlock");
    var b = $("<input>", {type: "submit"});
    b.click(getJson);
    d.append(b)
    return d
}

function textarea(prefix, name) {
    var d = $("<div/>");
    d.addClass("inputBlock");
    var label = $("<div/>").addClass("label");
    label.text(prefix);

    var textElement = $("<textarea/>", {
                        name: name,
                        rows: 5
    });
    textElement.addClass("inputElement");

    var errorContainer = $("<div/>", {name: "error:"+name});
    errorContainer.addClass("error");

    d.append(label);
    d.append(textElement);
    d.append(errorContainer);

    return d
}

function text(prefix, name) {
    var d = $("<div/>");
    d.addClass("inputBlock");

    var label = $("<div/>").addClass("label");

    label.text(prefix);
    var textElement = $("<input/>", {
                        type: "text",
                        name: name
    });
    textElement.addClass("inputElement");

    var errorContainer = $("<div/>", {name: "error:"+name});
    errorContainer.addClass("error");

    d.append(label);
    d.append(textElement);
    d.append(errorContainer);

    return d
}

function getJson() {
    var errorElements = $('[name^="error:"]').text(""); 
    var form = $("form")[0];
    var formdata = new FormData(form);
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/admin/user", false);
    xhr.send(formdata);
    var result = JSON.parse(xhr.responseText);
    var data = result.data;
    if (result.status == "ok") {
        alert("result witout errors")
    } else {
        for (err in data) {
            var errorDiv = $('[name="error:'+err+'"]');
            if (errorDiv.length == 0) {
                $('[name="error"]').text(data[err]);
            }
            errorDiv.text(data[err]);
        };
    };

    return false;
}
