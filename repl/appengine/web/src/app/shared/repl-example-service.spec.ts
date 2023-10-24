import { TestBed } from '@angular/core/testing';

import { Example, ReplExampleService } from './repl-example-service';

describe('ReplExampleService', () => {
  let service: ReplExampleService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ReplExampleService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should publish examples', () => {

    let recorded_example: Example = {description: "", request: {commands: []}};

    service.examplePosted$.subscribe({
      next: (example: Example) => {
        recorded_example = example;
      }
    });

    service.postExample({description: "description", request: {
      commands: [
        "%status"
      ]
    }});

    expect(recorded_example.description).toBe("description");
  });
});
