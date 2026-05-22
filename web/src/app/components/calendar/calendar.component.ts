import { Component, OnInit } from '@angular/core';



@Component({
  standalone: true,
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.css']
})
export class CalendarComponent implements OnInit {
  currentMonth: number=0;
  currentYear: number=0;
  daysInMonth: number[]=[];
  daysOfWeek: string[] = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb'];
  firstDayOfMonth: number=0;

  constructor() { }

  ngOnInit(): void {
    const date = new Date();
    this.currentMonth = date.getMonth();
    this.currentYear = date.getFullYear();
    this.generateCalendar();
  }

  generateCalendar(): void {
    const date = new Date(this.currentYear, this.currentMonth, 1);
    this.firstDayOfMonth = date.getDay();
    const daysInMonth = new Date(this.currentYear, this.currentMonth + 1, 0).getDate();
    
    this.daysInMonth = [];
    for (let i = 1; i <= daysInMonth; i++) {
      this.daysInMonth.push(i);
    }
  }

  changeMonth(direction: number): void {
    this.currentMonth += direction;
    if (this.currentMonth > 11) {
      this.currentMonth = 0;
      this.currentYear++;
    } else if (this.currentMonth < 0) {
      this.currentMonth = 11;
      this.currentYear--;
    }
    this.generateCalendar();
  }
}
