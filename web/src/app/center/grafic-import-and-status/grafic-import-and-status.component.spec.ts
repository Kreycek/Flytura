import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GraficImportAndStatusComponent } from './grafic-import-and-status.component';

describe('GraficImportAndStatusComponent', () => {
  let component: GraficImportAndStatusComponent;
  let fixture: ComponentFixture<GraficImportAndStatusComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GraficImportAndStatusComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GraficImportAndStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
