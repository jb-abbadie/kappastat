$(function() {
    $('#list_streams').hide();
    $('#list_duration').hide();
    $('#go_button').hide();
    $('#curve_chart').hide();
    $.getJSON("/api/following", listFollowing);
    var chart = initChart();
    updateChart(chart,[])
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
            for(var i in data) {
                result.push([new Date(data[i]['Start']), data[i]['Resub'], data[i]['Messages']]);
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
        explorer: {
            "axis":"horizontal",
        },
        animation: {
            "duration":1000,
            "easing":"inAndOut",
        },
        series : {
            0:{axis: 'Viewer'},
            1:{axis: 'Messages'},
        },
        vAxes: {
            0: {title: "Viewership"},
            1: {title: "Daylight"},
        },
    };
    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Time');
    data.addColumn('number', "Viewer");
    data.addColumn('number', "Chat Messages");
    data.addRows(inp);

    chart.draw(data, options);
}
