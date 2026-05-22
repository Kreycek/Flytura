import { Component, Input, input, Renderer2 } from '@angular/core';
import { CommonModule } from '@angular/common';


@Component({
  selector: 'app-spinner',
  imports: [CommonModule],
  templateUrl: './spinner.component.html',
  styleUrl: './spinner.component.css'
})
export class SpinnerComponent {


  @Input() isLoading = false;
  @Input() text = 'Processando…'; 

  constructor(private renderer: Renderer2) {}

  ngOnInit() {
    // Exemplo: ativa durante uma operação
    // this.setLoading(true);
    // setTimeout(() => this.setLoading(false), 2500);
  }

  setLoading(v: boolean) {
    this.isLoading = v;
    if (v) {
      this.renderer.addClass(document.body, 'overlay-active');
    } else {
      this.renderer.removeClass(document.body, 'overlay-active');
    }
  }

  ngOnDestroy() {
    this.renderer.removeClass(document.body, 'overlay-active');
  }

}
