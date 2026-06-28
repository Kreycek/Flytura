import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TotalGroupComponent } from './total-group.component';

describe('TotalGroupComponent', () => {
  let component: TotalGroupComponent;
  let fixture: ComponentFixture<TotalGroupComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TotalGroupComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TotalGroupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
