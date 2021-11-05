/* Just a simple timesheet program I wrote - Joseph Grimer Oct 2021
 *
 * an excuse to learn go and have a nice timesheet recorded and counted
 * at the same time
 */

package main

import (
    "errors"
    "fmt"
    "log"
    "os"
    "path"
    "time"
    "strconv"
)

const timesheet_dir = "/Users/joseph.grimer/timesheets"
const working_day = time.Duration(time.Hour * 7) + (time.Minute * 30)

func echoAppend(fh *os.File, strs ...string) {
    for _, str := range strs {
        _, err := fh.WriteString(str)
        fmt.Print(str)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func hourMin(d time.Duration) string {
    d = d.Round(time.Minute)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    return fmt.Sprintf("%d:%02d", h, m)
}

func main() {
    if _, err := os.Stat(timesheet_dir); errors.Is(err, os.ErrNotExist) {
        err = os.Mkdir(timesheet_dir, 0777);
        if err != nil {
            log.Fatal(err)
        }
    }

    job_started := time.Date(2021, 7, 25, 3, 0, 0, 0, time.Local)
    now := time.Now().Local()
    week_diff := strconv.Itoa(int(now.Sub(job_started).Hours()) / ( 24*7 ) + 1)

    // old working week - doesn't take into account potential holidays
    //working_week := working_day * time.Duration(6 - int(now.Weekday()))

    var on_the_clock = true
    var timecard [][]time.Time
    /*
    // for testing
    time_a := time.Date(2021, 11, 2, 3, 0, 0, 0, time.Local)
    time_b := time.Date(2021, 11, 3, 13, 0, 0, 0, time.Local)
    new_row := []time.Time{time_a, time_b}
    timecard = append(timecard, new_row)
    // */
    target_file := path.Join(timesheet_dir, week_diff)
    week_file, err := os.OpenFile(target_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
    if err != nil {
        log.Fatal(err)
    }
    var scan_string string
    echoAppend(week_file, "Week ", week_diff, "\n")
    echoAppend(week_file,"------------------------------------")
    for { // i.e. while true
        current_time := time.Now().Local()
        if on_the_clock {
            new_row := []time.Time{current_time}
            timecard = append(timecard, new_row)
        } else {
            timecard[len(timecard)-1] = append(timecard[len(timecard)-1], current_time)
        }
        
        first_day := true
        var last_period_start time.Time
        time_worked := time.Duration(0)
        working_week := time.Duration(0)  // in order to dynamically calculate time left to work, independant of holiday
        for _, row := range timecard {
            if first_day || last_period_start.Day() != row[0].Day(){
                echoAppend(week_file,"\n", row[0].Weekday().String(), " ")
                working_week += working_day
            }
            echoAppend(week_file,row[0].Format("3:04")," -> ")
            if len(row) > 1 {
                time_worked += row[1].Sub(row[0])
                echoAppend(week_file,row[1].Format("3:04"), " | ")
            } else {
                // predict
                start_time := row[0]
                time_left := working_week - time_worked
                if time_left < working_day {
                    if start_time.Hour() < 11 { // morning last day
                        echoAppend(week_file,"(", (start_time.Add(time_left).Add(time.Hour)).Format("3:04"), ")")
                    } else { // after lunch last day
                        echoAppend(week_file,"(", (start_time.Add(time_left)).Format("3:04"), ")")
                    }
                } else {
                    if start_time.Hour() < 11 { // morning normal day
                        echoAppend(week_file,"(", (start_time.Add(time.Hour * 8).Add(time.Minute*30)).Format("3:04"), ")")
                    } else { // after lunch normal day
                        echoAppend(week_file,"(", (last_period_start.Add(time.Hour * 8).Add(time.Minute*30)).Format("3:04"), ")")
                    }
                }
            }
            last_period_start = row[0]
            first_day = false
        }
        echoAppend(week_file,"\n------------------------------------\n")

        if on_the_clock {
            echoAppend(week_file,"You are working")
        } else {
            echoAppend(week_file,"You are leaving\n")
            time_left := working_week - time_worked
            echoAppend(week_file,"Worked ", hourMin(time_worked), " of ", hourMin(working_week),
            " leaving ", hourMin(time_left), "\n")
            echoAppend(week_file,"-------------------------------------")
        }

        // Taking input from user
        fmt.Scanln(&scan_string)
        on_the_clock = !on_the_clock
    }
}
