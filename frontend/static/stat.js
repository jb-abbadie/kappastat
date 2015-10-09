$(function() {
    $('#list_streams').hide();
    $('#list_duration').hide();
    $('#go_button').hide();
    $.getJSON("/api/following", listFollowing);
    var chart = initChart();
    updateChart(chart,[["Time","Viewers"],[0,0],[1,0]])
    console.log("initialized");

    $('#go_button').click( function() {
        var selected = $('#list_streams option:selected').attr("value");
        var duration = $('#list_duration option:selected').attr("time");
        if ( duration === undefined) {
            duration = 15
        }
        if ( selected !== undefined) {
        $.getJSON("/api/stat/" + selected ,{"duration":duration}, function(data) {

            var result = [];
            result.push(addTableHeader());
            console.log(result)
            for(var i in data) {
                result.push(filterResult(data[i]));
            }
            updateChart(chart, result);
        });
        }
    });
});

function initChart() {
    var chart = new google.visualization.LineChart(document.getElementById('curve_chart'));
    return chart;
}

function updateChart(chart, inp) {
    var options = {
        title: 'This is a chart',
        fontName: "Raleway",
        fontSize: 15,
        explorer: {
            "axis":"horizontal",
        },
        animation: {
            "duration":1000,
            "easing":"inAndOut",
        },
    };
    console.log(data)
    var data = new google.visualization.arrayToDataTable(inp);

    chart.draw(data, options);
}

function addTableHeader() {
    var ret = [];
    ret.push("Time");
    if ($("#cb-views").is(':checked()')) {
        ret.push("Viewer");
    }
    if ($("#cb-chat").is(':checked()')) {
        ret.push("Messages");
    }
    if ($("#cb-sub").is(':checked()')) {
        ret.push("Newsub");
    }
    if ($("#cb-resub").is(':checked()')) {
        ret.push("Resub");
    }
    return ret;
}

function filterResult(data) {
    var ret = [];
    ret.push(new Date(data["Start"]));
    if ($("#cb-views").is(':checked()')) {
        ret.push(data["Viewer"]);
    }
    if ($("#cb-chat").is(':checked()')) {
        ret.push(data["Messages"]);
    }
    if ($("#cb-sub").is(':checked()')) {
        ret.push(data["Newsub"]);
    }
    if ($("#cb-resub").is(':checked()')) {
        ret.push(data["Resub"]);
    }
    return ret;
}
