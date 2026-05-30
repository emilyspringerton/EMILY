# Emily: Iterative Code Execution Loop with Acceptance Criteria

## Overview

Emily can spawn and manage iterative Claude Code execution loops—continuously running code, testing against acceptance criteria, and iterating until all criteria pass. This is her core mechanism for autonomous problem-solving and self-improvement.

```
┌─ Emily: "Optimize deployment script until it runs in < 30 seconds"
│
├─ Define Acceptance Criteria
│  ├─ Execution time < 30 seconds
│  ├─ All health checks pass
│  └─ No errors in logs
│
├─ Loop: Run Claude Code
│  ├─ Iteration 1: Run script → measure → fails (45s)
│  │  └─ Analyze: "Pre-checks are parallelizable"
│  │
│  ├─ Iteration 2: Refactor → measure → fails (32s, closer)
│  │  └─ Analyze: "API calls can be batched"
│  │
│  ├─ Iteration 3: Batch APIs → measure → passes! (28s)
│  │  └─ Success! Criteria met.
│  │
│  └─ Log: Improvement captured, pattern stored
│
└─ Result: Optimized script, lessons learned
```

---

## 1. ACCEPTANCE CRITERIA SPECIFICATION

### 1.1 Criteria Types

**Type A: Measurable Performance**
```yaml
Criteria:
  - name: "execution_time"
    metric: "seconds"
    threshold: "<"
    target: 30
    measurement: "time script execution with `time` command"
    
  - name: "cost"
    metric: "dollars"
    threshold: "<"
    target: 0.05
    measurement: "sum API call costs logged in results"
    
  - name: "error_rate"
    metric: "percent"
    threshold: "<"
    target: 1
    measurement: "failed operations / total operations"
```

**Type B: Functional Correctness**
```yaml
Criteria:
  - name: "test_suite_pass"
    type: "boolean"
    requirement: true
    measurement: "run pytest, check exit code"
    
  - name: "schema_valid"
    type: "boolean"
    requirement: true
    measurement: "validate generated SQL against schema"
```

**Type C: Qualitative/Heuristic**
```yaml
Criteria:
  - name: "code_quality"
    type: "heuristic"
    rubric: |
      - Is there docstring on all functions?
      - Is variable naming clear?
      - Are there obvious optimizations missed?
      - Could a junior engineer understand this?
    measurement: "Claude Code Artifact review"
```

**Type D: Safety/Compliance**
```yaml
Criteria:
  - name: "no_secrets_in_logs"
    type: "boolean"
    requirement: true
    measurement: "grep for AWS_KEY, DB_PASSWORD, API_TOKEN patterns"
    
  - name: "rollback_plan_exists"
    type: "boolean"
    requirement: true
    measurement: "Check for rollback script in deployment artifact"
```

### 1.2 Acceptance Criteria Document

```yaml
improvement_task: "Optimize financial transaction processing"
acceptance_criteria:
  - criterion: "latency"
    metric: "p99_latency_ms"
    current: 450
    target: 300
    threshold: "<="
    weight: "high"
    
  - criterion: "throughput"
    metric: "transactions_per_second"
    current: 5000
    target: 7000
    threshold: ">="
    weight: "high"
    
  - criterion: "consistency"
    metric: "audit_log_complete"
    current: true
    target: true
    threshold: "=="
    weight: "critical"
    
  - criterion: "safety"
    metric: "no_lost_transactions"
    current: true
    target: true
    threshold: "=="
    weight: "critical"
    
  - criterion: "code_coverage"
    metric: "percent"
    current: 82
    target: 90
    threshold: ">="
    weight: "medium"

all_criteria_must_pass: true  # vs. weighted scoring
max_iterations: 20
max_duration_minutes: 120
timeout_per_iteration_minutes: 10
```

---

## 2. THE ITERATIVE LOOP PATTERN

### 2.1 Loop Architecture

```python
# Pseudocode of Emily's iterative loop engine

class IterativeImprovement:
    def __init__(self, task_description, acceptance_criteria, max_iterations=20):
        self.task = task_description
        self.criteria = acceptance_criteria
        self.max_iterations = max_iterations
        self.iteration = 0
        self.execution_trace = []
        
    def run_loop(self):
        while self.iteration < self.max_iterations:
            self.iteration += 1
            
            # Step 1: Generate code (or get user input for first iteration)
            code = self.get_code_for_iteration()
            
            # Step 2: Execute code in Claude environment
            result = self.execute_code(code)
            
            # Step 3: Measure against acceptance criteria
            measurements = self.measure_against_criteria(result)
            
            # Step 4: Check if all criteria pass
            if self.all_criteria_pass(measurements):
                return self.success(measurements)
            
            # Step 5: If not, analyze failure and prepare next iteration
            analysis = self.analyze_failure(measurements, result)
            self.execution_trace.append({
                'iteration': self.iteration,
                'code': code,
                'measurements': measurements,
                'analysis': analysis,
                'next_focus': self.identify_next_focus(analysis)
            })
            
        # Loop exhausted without success
        return self.partial_success(self.execution_trace)
```

### 2.2 Iteration Flow (Detailed)

```
ITERATION 1
│
├─ Generate Code
│  └─ Emily: "Write Python script to benchmark transaction processor"
│     ├─ Call Claude API (claude-sonnet-4) with:
│     │  - Task description
│     │  - Acceptance criteria
│     │  - Previous iterations (if any)
│     │  └─ System prompt optimized for iterative improvement
│     └─ Returns: Code artifact
│
├─ Execute Code (in Claude sandbox)
│  ├─ Run the generated script
│  ├─ Capture stdout, stderr, return code
│  ├─ Set timeout (per iteration)
│  └─ Return: execution result + logs
│
├─ Measure Against Criteria
│  ├─ p99_latency: Parse logs → 450ms ❌ (target: 300ms)
│  ├─ throughput: Parse metrics → 5000 TPS ❌ (target: 7000 TPS)
│  ├─ consistency: Audit check → ✅ (target: pass)
│  ├─ safety: Transaction count → ✅ (target: pass)
│  └─ code_coverage: pytest → 82% ❌ (target: 90%)
│
├─ Evaluate: All criteria pass?
│  └─ NO → Continue to Step 5
│
├─ Analyze Failure
│  ├─ Which criteria failed? (3/5)
│  ├─ By how much? (50ms off on latency, 2000 TPS off on throughput)
│  ├─ Are failures related? (yes—both are performance issues)
│  ├─ What was tried? (baseline implementation)
│  ├─ What might work? (batch API calls, add caching, parallelize)
│  ├─ Why wasn't it done? (first iteration is always baseline)
│  └─ What to try next? (batch API calls)
│
└─ Prepare Next Iteration
   ├─ Add to execution trace
   ├─ Generate new system prompt for next iteration
   │  ├─ Include: failed measurements + analysis
   │  ├─ Include: what was tried + what didn't work
   │  ├─ Include: the three optimization hypotheses
   │  └─ Instruct: "Focus on batching API calls next"
   └─ Loop back to "Generate Code"

ITERATION 2
│
├─ Generate Code (with context from Iteration 1)
│  └─ Claude: "Refactor to batch API calls, maintain safety/consistency"
│
├─ Execute → Measure
│  ├─ p99_latency: 380ms ✅ improved but still ❌ (target: 300ms)
│  ├─ throughput: 5800 TPS ❌ (target: 7000 TPS)
│  ├─ consistency: ✅
│  ├─ safety: ✅
│  └─ code_coverage: 85% ❌ (target: 90%)
│
├─ Analyze: Made progress! Latency improved 70ms. But throughput lagging.
│
└─ Next Focus: "Parallelize transaction processing"

ITERATION 3
│
├─ Generate Code (with history of Iterations 1-2)
│  └─ Claude: "Add parallel processing for independent transactions"
│
├─ Execute → Measure
│  ├─ p99_latency: 295ms ✅ (target: 300ms) ✅✅
│  ├─ throughput: 7200 TPS ✅ (target: 7000 TPS) ✅✅
│  ├─ consistency: ✅
│  ├─ safety: ✅
│  └─ code_coverage: 92% ✅ (target: 90%) ✅✅
│
├─ Evaluate: ALL CRITERIA PASS! ✅✅✅
│
└─ Success!
   ├─ Total iterations: 3
   ├─ Final code committed
   ├─ Lessons learned captured:
   │  ├─ Pattern: "Batch operations → Big latency wins"
   │  ├─ Pattern: "Parallelization needs careful audit logging"
   │  ├─ Pattern: "Start with baseline, then optimize"
   │  └─ Anti-pattern: "Over-optimize before measuring"
   ├─ Improvement metrics logged
   └─ Added to decision heuristics for future tasks
```

---

## 3. SYSTEM PROMPTS FOR ITERATION

### 3.1 First Iteration Prompt

```
You are Emily, an AI chief of staff agent.

Task: Optimize financial transaction processing
Acceptance Criteria:
  - p99_latency < 300ms (currently 450ms)
  - throughput > 7000 TPS (currently 5000 TPS)
  - audit log consistency = pass
  - no lost transactions = pass
  - code coverage > 90% (currently 82%)

Strategy for Iteration 1:
- Write a baseline implementation that is CORRECT first
- Instrument it heavily with metrics and logging
- We'll optimize in later iterations
- Include tests for all critical paths
- Focus on clarity over speed in this iteration

Return ONLY a working Python script that:
1. Processes a simulated batch of transactions
2. Logs metrics (latency, throughput, errors)
3. Validates consistency
4. Includes test cases

Ready? Generate the code.
```

### 3.2 Subsequent Iteration Prompt

```
You are Emily, an AI chief of staff agent, iterating on a task.

Task: Optimize financial transaction processing

Previous Results:
Iteration 1:
  - p99_latency: 450ms ❌ (target: 300ms, gap: -150ms)
  - throughput: 5000 TPS ❌ (target: 7000 TPS, gap: -2000)
  - consistency: ✅
  - safety: ✅
  - code_coverage: 82% ❌ (target: 90%, gap: -8%)
  
Analysis:
- Performance is the limiting factor (both latency & throughput)
- Consistency and safety are solid
- Code is readable but has room for optimization
- Root cause analysis suggests:
  - Sequential processing of transactions blocks throughput
  - API calls to ledger service are slow
  - Could batch operations

What Worked in Iteration 1:
- Clear audit logging architecture
- Transaction consistency checks
- Error handling

What Didn't Work:
- Sequential processing creates bottleneck
- No caching of ledger lookups
- Tests are thorough but slow

Hypothesis for Iteration 2:
- Batch API calls to ledger service
- Add memoization for frequent lookups
- These should improve both latency and throughput

Please refactor the code to:
1. Batch ledger API calls (group 10-20 transactions per batch)
2. Add a simple in-memory cache for ledger balances
3. Maintain all safety and consistency guarantees
4. Improve test coverage to >85% if possible

Focus on: PERFORMANCE OPTIMIZATION (latency & throughput)
Do NOT compromise: Safety, consistency, or audit trails

Return the refactored code only.
```

---

## 4. MEASUREMENT & EVALUATION

### 4.1 Measurement Execution

```python
def measure_against_criteria(execution_result):
    """
    Measure code execution against acceptance criteria
    """
    measurements = {}
    
    # Parse metrics from execution logs
    logs = execution_result.stdout + execution_result.stderr
    
    # Criterion 1: Latency
    latency_match = re.search(r'p99_latency: (\d+)ms', logs)
    if latency_match:
        measurements['p99_latency'] = {
            'value': int(latency_match.group(1)),
            'target': 300,
            'threshold': '<',
            'passes': int(latency_match.group(1)) < 300
        }
    
    # Criterion 2: Throughput
    tps_match = re.search(r'throughput: (\d+) TPS', logs)
    if tps_match:
        measurements['throughput'] = {
            'value': int(tps_match.group(1)),
            'target': 7000,
            'threshold': '>=',
            'passes': int(tps_match.group(1)) >= 7000
        }
    
    # Criterion 3: Consistency check
    consistency_check = 'consistency: PASS' in logs
    measurements['consistency'] = {
        'value': 'pass' if consistency_check else 'fail',
        'target': 'pass',
        'threshold': '==',
        'passes': consistency_check
    }
    
    # ...more criteria...
    
    return measurements
```

### 4.2 Evaluation Logic

```python
def all_criteria_pass(measurements):
    """Check if all criteria are met"""
    for criterion_name, measurement in measurements.items():
        if not measurement['passes']:
            return False
    return True

def partial_success_score(measurements):
    """
    If we exhaust iterations without perfect success,
    score how close we got (useful for learning)
    """
    criteria_passed = sum(1 for m in measurements.values() if m['passes'])
    total_criteria = len(measurements)
    
    gap_score = 0
    for measurement in measurements.values():
        if not measurement['passes']:
            gap = abs(measurement['value'] - measurement['target'])
            gap_percent = (gap / measurement['target']) * 100
            gap_score += gap_percent
    
    return {
        'criteria_passed': criteria_passed,
        'total_criteria': total_criteria,
        'total_gap_percent': gap_score,
        'partial_success': True
    }
```

---

## 5. FAILURE ANALYSIS & NEXT ITERATION PLANNING

### 5.1 Failure Analysis

```python
def analyze_failure(measurements, code, execution_result):
    """
    Analyze why criteria weren't met and suggest next steps
    """
    analysis = {
        'iteration': self.iteration,
        'timestamp': now(),
        'criteria_failed': [],
        'criteria_passed': [],
        'failure_patterns': [],
        'hypotheses': [],
        'root_causes': [],
        'next_focus_areas': []
    }
    
    for name, m in measurements.items():
        if m['passes']:
            analysis['criteria_passed'].append(name)
        else:
            analysis['criteria_failed'].append({
                'name': name,
                'current': m['value'],
                'target': m['target'],
                'gap': m['value'] - m['target'],
                'gap_percent': ((m['value'] - m['target']) / m['target']) * 100
            })
    
    # Pattern detection
    if 'p99_latency' in [f['name'] for f in analysis['criteria_failed']] and \
       'throughput' in [f['name'] for f in analysis['criteria_failed']]:
        analysis['failure_patterns'].append(
            'PERFORMANCE_BOTTLENECK: Both latency and throughput failing suggests sequential processing issue'
        )
    
    # Hypotheses based on code inspection
    if 'for transaction in transactions:' in code:
        analysis['hypotheses'].append(
            'HYPOTHESIS: Sequential loop over transactions. Try batching or parallelizing.'
        )
    
    # Root cause analysis from logs
    if 'API call' in execution_result.stderr and 'timeout' in execution_result.stderr.lower():
        analysis['root_causes'].append(
            'ROOT_CAUSE: API timeouts detected. Suggest batching to reduce call count.'
        )
    
    # Next focus areas (ranked by impact)
    analysis['next_focus_areas'] = [
        {
            'focus': 'Parallelize transaction processing',
            'expected_impact': 'Could improve throughput by 50% and latency by 40%',
            'effort': 'medium',
            'risk': 'low (consistency guarantees must be maintained)',
            'priority': 1
        },
        {
            'focus': 'Batch API calls',
            'expected_impact': 'Could improve latency by 30%',
            'effort': 'low',
            'risk': 'very_low',
            'priority': 2
        },
        {
            'focus': 'Add caching for frequent lookups',
            'expected_impact': 'Could improve throughput by 20%',
            'effort': 'low',
            'risk': 'low (must invalidate cache correctly)',
            'priority': 3
        }
    ]
    
    return analysis
```

### 5.2 Next Iteration Prompt Generation

```python
def generate_next_iteration_prompt(task, criteria, execution_trace):
    """
    Generate the system prompt for the next Claude Code iteration
    """
    prompt = f"""
You are Emily, an AI chief of staff agent, iterating on a task.

TASK: {task['description']}

ACCEPTANCE CRITERIA (All must pass):
{format_criteria(criteria)}

PREVIOUS ITERATIONS RESULTS:
"""
    
    # Add execution history
    for i, trace in enumerate(execution_trace, 1):
        prompt += f"""
Iteration {i}:
  Measurements:
{format_measurements(trace['measurements'])}
  
  Analysis: {trace['analysis']['failure_patterns']}
  
  What was tried: {trace['code_summary']}
"""
    
    # Add focused guidance
    if execution_trace:
        last_analysis = execution_trace[-1]['analysis']
        prompt += f"""

NEXT ITERATION FOCUS (Iteration {len(execution_trace) + 1}):
Priority improvements to try:
"""
        for focus in last_analysis['next_focus_areas']:
            prompt += f"  - {focus['focus']} (Expected impact: {focus['expected_impact']})\n"
    
    prompt += """

CONSTRAINTS:
- Do NOT compromise safety or consistency
- Maintain audit logging
- Keep code readable
- Include tests

Return ONLY the refactored code.
"""
    
    return prompt
```

---

## 6. TERMINATION CONDITIONS

### 6.1 Success Termination

```
All acceptance criteria pass → COMPLETE
├─ Log: Success metrics
├─ Commit: Final code to repo
├─ Learn: Extract patterns to knowledge base
└─ Return: Execution trace + lessons learned
```

### 6.2 Timeout Termination

```
Max iterations exceeded (e.g., 20 attempts) → PARTIAL SUCCESS
├─ Evaluate: How close did we get?
├─ Log: Best-effort result + why we stopped
├─ Decide: 
│  ├─ If gap is small: Accept partial solution
│  ├─ If gap is large: Escalate to human
│  └─ Learn: Why was the problem hard?
└─ Return: Execution trace + analysis
```

### 6.3 Cost/Duration Termination

```
Time or token budget exceeded → HALT WITH ANALYSIS
├─ This task is expensive or complex
├─ Escalate to Emily + human with:
│  ├─ What was accomplished
│  ├─ Where we got stuck
│  └─ Recommended next steps
└─ Learning: Add to "hard problems" backlog
```

### 6.4 Safety Termination

```
Detected unsafe operation (e.g., deleting prod database) → IMMEDIATE HALT
├─ Do not execute
├─ Escalate to human immediately
├─ Log: What code was trying to do
└─ Question: Can we constrain this task better?
```

---

## 7. EXECUTION TRACE & LEARNING

### 7.1 Execution Trace Format

```json
{
  "task_id": "improve-transaction-processing-2025-05-30",
  "task_description": "Optimize financial transaction processing",
  "started": "2025-05-30T14:00:00Z",
  "completed": "2025-05-30T14:28:45Z",
  "total_iterations": 3,
  "total_duration_seconds": 1725,
  "result": "SUCCESS",
  
  "acceptance_criteria": {
    "p99_latency": {"target": 300, "final_value": 295, "passes": true},
    "throughput": {"target": 7000, "final_value": 7200, "passes": true},
    "consistency": {"target": true, "final_value": true, "passes": true},
    "safety": {"target": true, "final_value": true, "passes": true},
    "code_coverage": {"target": 90, "final_value": 92, "passes": true}
  },
  
  "iterations": [
    {
      "iteration_number": 1,
      "duration_seconds": 45,
      "tokens_used": 1250,
      "code_summary": "Sequential transaction processing with detailed logging",
      "measurements": {
        "p99_latency": 450,
        "throughput": 5000,
        "consistency": true,
        "safety": true,
        "code_coverage": 82
      },
      "analysis": {
        "patterns": ["Sequential processing bottleneck"],
        "next_focus": "Batch API calls and parallelize"
      }
    },
    {
      "iteration_number": 2,
      "duration_seconds": 52,
      "tokens_used": 1400,
      "code_summary": "Batch API calls, added memoization",
      "measurements": {
        "p99_latency": 380,
        "throughput": 5800,
        "consistency": true,
        "safety": true,
        "code_coverage": 85
      },
      "analysis": {
        "patterns": ["Good progress on latency but throughput still limited"],
        "next_focus": "Add parallel processing"
      }
    },
    {
      "iteration_number": 3,
      "duration_seconds": 68,
      "tokens_used": 1600,
      "code_summary": "Parallel transaction processing with safety guarantees",
      "measurements": {
        "p99_latency": 295,
        "throughput": 7200,
        "consistency": true,
        "safety": true,
        "code_coverage": 92
      },
      "analysis": {
        "patterns": ["All criteria met!"],
        "lessons_learned": [
          "Pattern: Parallelization + batching = 40% latency improvement",
          "Pattern: Audit logging must happen at boundaries, not per-transaction",
          "Anti-pattern: Premature optimization before measuring"
        ]
      }
    }
  ],
  
  "lessons_learned": [
    "Batch operations before optimizing individual transaction speed",
    "Parallelization requires careful consistency checking",
    "Instrument first, optimize second"
  ],
  
  "patterns_to_remember": [
    {
      "name": "Performance bottleneck detection",
      "pattern": "When both latency and throughput fail, look for sequential loops",
      "solution": "Batch or parallelize operations"
    },
    {
      "name": "Audit logging at scale",
      "pattern": "Per-transaction logging doesn't scale; batch log writes",
      "solution": "Accumulate changes in memory, flush periodically"
    }
  ],
  
  "added_to_knowledge_base": [
    "financial-transaction-optimization/batch-api-pattern.md",
    "financial-transaction-optimization/parallel-safety-guardrails.md"
  ]
}
```

### 7.2 Pattern Extraction & Storage

```
After successful completion, Emily extracts:

1. WHAT WORKED
   └─ Patterns that led to success
      └─ Stored in: knowledge-base/optimization-patterns.md

2. WHAT DIDN'T
   └─ Anti-patterns & dead ends
      └─ Stored in: knowledge-base/optimization-anti-patterns.md

3. DECISION LEARNINGS
   └─ When to batch vs. parallelize? When to cache?
      └─ Stored in: decision-heuristics/v2.3.md

4. DOMAIN KNOWLEDGE
   └─ How financial transactions work, safety requirements
      └─ Stored in: knowledge-base/financial-domain.md

5. TOOL EFFECTIVENESS
   └─ API batching reduced latency by 70ms; parallelization by 80ms
      └─ Stored in: knowledge-base/technique-effectiveness.md
```

---

## 8. APPLYING THIS TO SELF-IMPROVEMENT

Emily can use this loop for her own recursive improvement:

```
Goal: Improve decision accuracy from 87% to 95%

Acceptance Criteria:
  - Decision accuracy > 95%
  - Decision time < 2 minutes per decision
  - False positive rate (escalations that weren't needed) < 5%
  - Team satisfaction with decisions > 8.5/10

Iteration 1:
├─ Baseline decision heuristics v1.0
├─ Measure accuracy on historical decisions
└─ Result: 87% accuracy (baseline)

Iteration 2:
├─ Add pattern matching for common scenarios
├─ Update decision tree for finance operations
└─ Result: 90% accuracy (improvement!)

Iteration 3:
├─ Integrate feedback from specialist agents
├─ Add confidence scoring to decisions
└─ Result: 93% accuracy (getting closer)

Iteration 4:
├─ Build meta-rules: When to escalate?
├─ Test on edge cases
└─ Result: 95% accuracy ✅ SUCCESS

Final Result:
├─ New heuristics committed: decision-heuristics/v2.4.md
├─ Lessons: Specialist feedback + pattern matching = biggest wins
├─ Next focus: Reduce decision time (currently 3m, target 2m)
└─ Ready for deployment
```

---

## 9. EXAMPLE: DEPLOYMENT OPTIMIZATION TASK

Let's say Emily is tasked with: **"Make deployments 50% faster while maintaining safety"**

### Iteration 1: Baseline

**Code Generated:**
```python
# Current deployment process (baseline measurement)
def deploy_service(service_name):
    start = time.time()
    
    # Phase 1: Pre-checks (currently ~120 seconds)
    run_tests()                    # 60s
    lint_code()                    # 15s
    security_scan()                # 30s
    health_check_staging()         # 15s
    
    # Phase 2: Deployment (currently ~180 seconds)
    build_docker_image()           # 80s
    push_to_registry()             # 20s
    update_kubernetes()            # 45s
    run_smoke_tests()              # 35s
    
    # Phase 3: Validation (currently ~60 seconds)
    verify_metrics()               # 30s
    check_logs_for_errors()        # 30s
    
    total = time.time() - start
    return {
        'duration': total,
        'success': total < 400  # Current: ~360s
    }
```

**Iteration 1 Measurement:** 360 seconds ❌ (Target: 180 seconds)

### Iteration 2: Parallelize Pre-checks

**Analysis of Iteration 1:**
- Pre-checks (120s) are independent → can run in parallel
- Currently: test → lint → security → health check (sequential)
- Hypothesis: Run in parallel, reduce to ~60s

**Code Generated:**
```python
def deploy_service(service_name):
    start = time.time()
    
    # Phase 1: Pre-checks (parallelized)
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = [
            executor.submit(run_tests),           # 60s
            executor.submit(lint_code),           # 15s
            executor.submit(security_scan),       # 30s
            executor.submit(health_check_staging) # 15s
        ]
        concurrent.futures.wait(futures)
        # Total: ~60s (max of all, not sum)
    
    # Phase 2: Deployment (unchanged for now)
    build_docker_image()           # 80s
    push_to_registry()             # 20s
    update_kubernetes()            # 45s
    run_smoke_tests()              # 35s
    
    # Phase 3: Validation (unchanged)
    verify_metrics()               # 30s
    check_logs_for_errors()        # 30s
    
    return {'duration': time.time() - start, ...}
```

**Iteration 2 Measurement:** 245 seconds ❌ (Target: 180 seconds, but 32% faster!)

### Iteration 3: Parallelize Deployment

**Analysis of Iteration 2:**
- Phase 2 (180s) is mostly blocking: Docker build, push, K8s update
- Hypothesis: Build Docker image while running final pre-checks
- Hypothesis: Push to registry while K8s update in progress

**Code Generated:**
```python
def deploy_service(service_name):
    start = time.time()
    
    # Pre-checks (parallel, ~60s)
    with ThreadPoolExecutor() as executor:
        pre_check_futures = [
            executor.submit(run_tests),
            executor.submit(lint_code),
            executor.submit(security_scan),
            executor.submit(health_check_staging)
        ]
        concurrent.futures.wait(pre_check_futures)
    
    # Deploy (optimized pipeline)
    # Start Docker build immediately
    build_future = executor.submit(build_docker_image)  # 80s
    
    # While Docker builds, can do other work
    push_registry_future = executor.submit(
        lambda: build_future.result() and push_to_registry()
    )  # Waits for build, then pushes (~20s after build done)
    
    # Update K8s in parallel with docker build+push
    update_k8s_future = executor.submit(
        lambda: push_registry_future.result() and update_kubernetes()
    )  # Waits for push, then updates (~45s)
    
    # Smoke tests after K8s update
    update_k8s_future.result()
    run_smoke_tests()  # 35s (after update)
    
    # Validation (parallel)
    with ThreadPoolExecutor() as executor:
        validation_futures = [
            executor.submit(verify_metrics),
            executor.submit(check_logs_for_errors)
        ]
        concurrent.futures.wait(validation_futures)
    
    return {'duration': time.time() - start, ...}
```

**Iteration 3 Measurement:** 190 seconds ✅ (Target: 180 seconds, very close!)

### Iteration 4: Fine-tuning

**Analysis of Iteration 3:**
- We're 10 seconds over. Small optimizations remaining.
- Could optimize image layer caching (saves ~10s if layers unchanged)
- Could parallelize smoke tests with validation

**Code Generated:**
```python
def deploy_service(service_name):
    # ...same pre-checks...
    
    # Deployment pipeline with optimization
    build_future = executor.submit(
        build_docker_image_with_cache_optimization()
    )  # 75s (cache hit reduces to 75s)
    
    # ...rest of pipeline...
    
    # Smoke tests + validation in parallel
    smoke_tests_future = executor.submit(run_smoke_tests)
    validation_futures = [
        executor.submit(verify_metrics),
        executor.submit(check_logs_for_errors)
    ]
    concurrent.futures.wait([smoke_tests_future] + validation_futures)
```

**Iteration 4 Measurement:** 175 seconds ✅✅ (Target: 180 seconds, SUCCESS!)

**Acceptance Criteria:**
- Deployment speed: 175s vs. 360s = **51% faster** ✅
- Safety maintained: All safety checks still run ✅
- Error detection: Same reliability ✅

**Lessons Learned:**
- Pattern: "Parallelize independent operations → 30-50% speedup"
- Pattern: "Pipeline sequential dependencies → avoid unnecessary waits"
- Anti-pattern: "Don't optimize before profiling"
- Storage: Added to `decision-heuristics/deployment-optimization.md`

---

## 10. EMILY'S LOOP COMMAND INTERFACE

How you ask Emily to run an improvement loop:

```
User: "Emily, optimize the deployment process. It should complete in under 180 seconds 
while maintaining all safety checks. Try iteratively for up to 20 attempts."

Emily: "Understood. Starting iterative improvement loop.
  Task: Deployment optimization
  Target: < 180 seconds
  Safety constraints: All checks maintained
  Max iterations: 20
  
  Starting iteration 1..."

[Loop runs, Emily reports after each iteration]

Emily: "Iteration 1 complete: 360s (baseline)"
Emily: "Iteration 2 complete: 245s (parallelized pre-checks)"
Emily: "Iteration 3 complete: 190s (pipeline optimization)"
Emily: "Iteration 4 complete: 175s (cache tuning)"

Emily: "✅ SUCCESS! All acceptance criteria met in 4 iterations.
  - Final speed: 175s (51% improvement)
  - Safety: All checks maintained
  - Lessons learned: 3 patterns extracted
  - Result committed to: repo/deploy/optimized_deploy.py
  - Execution trace: logs/improvement-2025-05-30.json"
```

---

## 11. TECHNICAL IMPLEMENTATION

### 11.1 Claude API Integration

Emily makes iterative calls to Claude API:

```python
import anthropic

class IterativeLoopEngine:
    def __init__(self):
        self.client = anthropic.Anthropic()
        self.model = "claude-opus-4"  # Latest model
        self.max_tokens = 4000  # Code artifacts
        
    def run_iteration(self, iteration_num, task, criteria, history):
        """Run a single iteration of the loop"""
        
        # Build system prompt
        system = self.build_system_prompt(task, criteria, history, iteration_num)
        
        # Build user message
        user_msg = self.build_user_message(iteration_num, task, history)
        
        # Call Claude
        response = self.client.messages.create(
            model=self.model,
            max_tokens=self.max_tokens,
            system=system,
            messages=[
                {"role": "user", "content": user_msg}
            ]
        )
        
        # Extract code from response
        code = self.extract_code(response.content)
        
        # Execute code in sandbox
        exec_result = self.execute_code_safely(code)
        
        # Measure results
        measurements = self.measure(exec_result, criteria)
        
        return {
            'code': code,
            'execution_result': exec_result,
            'measurements': measurements,
            'all_pass': all(m['passes'] for m in measurements.values())
        }
```

### 11.2 Execution Sandbox

```python
def execute_code_safely(code, timeout_seconds=60):
    """
    Execute generated code in a sandbox with:
    - Timeout protection
    - Memory limits
    - File system constraints
    - Network access controls
    """
    import subprocess
    import tempfile
    
    with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
        f.write(code)
        f.flush()
        
        try:
            result = subprocess.run(
                ['/usr/bin/timeout', str(timeout_seconds), 'python', f.name],
                capture_output=True,
                text=True,
                timeout=timeout_seconds + 5
            )
            
            return {
                'stdout': result.stdout,
                'stderr': result.stderr,
                'return_code': result.returncode,
                'success': result.returncode == 0
            }
        except subprocess.TimeoutExpired:
            return {
                'stdout': '',
                'stderr': f'Timeout after {timeout_seconds}s',
                'return_code': -1,
                'success': False
            }
```

---

## 12. LIMITATIONS & SAFETY CONSIDERATIONS

### 12.1 Infinite Loop Prevention

```yaml
Safeguards:
  - Max iterations: Hard limit (usually 20)
  - Timeout per iteration: Kill if > 10 minutes
  - Total duration limit: Stop if > 2 hours
  - Cost tracking: Alert if > $10 spent on one task
  - Divergence detection: If measurements getting worse, halt & escalate
```

### 12.2 When to Escalate Instead of Loop

```
❌ Don't loop if:
   - Problem is ambiguous or under-specified
   - Acceptance criteria are contradictory
   - Requires human judgment or domain knowledge you lack
   - Would affect production systems without human approval
   - Cost/time budget is too tight for iteration

✅ Do loop if:
   - Problem is well-defined
   - Criteria are measurable and achievable
   - Failure is low-risk (dev/staging only)
   - There's time/budget for iteration
   - Task is tactical (optimization, bug fix, refactoring)
```

### 12.3 Safety Constraints in the Loop

Hard rules Emily cannot violate, even in a loop:
```
- Never delete data without explicit approval
- Never bypass security controls
- Never expose secrets in generated code
- Never make prod changes without safety gates
- Always maintain audit trails
- Always include rollback plans
```

---

## 13. SUCCESS METRICS FOR THE LOOP

Emily's iterative loop is working well when:

✅ **Convergence:** Measurements improve each iteration toward criteria
✅ **Efficiency:** Reaches success in 3-6 iterations (not 15+)
✅ **Safety:** No iteration produces unsafe or broken code
✅ **Learning:** Patterns are extracted and reused
✅ **Transparency:** Execution trace is complete and auditable
✅ **Team Trust:** Team understands and agrees with optimizations

---

## Next Steps

1. **Define your first improvement task:** What problem should Emily optimize iteratively?
2. **Write detailed acceptance criteria:** Make them measurable
3. **Set iteration limits:** Max attempts, timeout, cost budget
4. **Test the loop:** Run 1-2 iterations manually to validate approach
5. **Integrate learnings:** Feed results back into Emily's decision heuristics
